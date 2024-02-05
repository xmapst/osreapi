package worker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/xmapst/osreapi/internal/dag"
	"github.com/xmapst/osreapi/internal/exec"
	"github.com/xmapst/osreapi/internal/logx"
	_ "github.com/xmapst/osreapi/internal/plugins"
	"github.com/xmapst/osreapi/internal/server/config"
	"github.com/xmapst/osreapi/internal/storage"
	"github.com/xmapst/osreapi/internal/storage/types"
	"github.com/xmapst/osreapi/internal/utils"
)

type TaskStep struct {
	ID             string
	CommandType    string
	CommandContent string
	EnvVars        []string
	DependsOn      []string
	Timeout        time.Duration

	state *types.TaskStepState
}

type Task struct {
	scriptDir string
	workspace string
	graph     *dag.Graph
	state     *types.TaskState

	ID       string
	Timeout  time.Duration
	EnvVars  []string
	MetaData types.MetaData
	Steps    []*TaskStep
}

func Submit(task *Task) error {
	atomic.AddInt64(&taskTotal, 1)
	task.workspace = filepath.Join(config.App.WorkSpace, task.ID)
	task.scriptDir = filepath.Join(config.App.ScriptDir, task.ID)
	task.state = &types.TaskState{
		ID:       task.ID,
		State:    exec.Pending,
		Count:    int64(len(task.Steps)),
		EnvVars:  task.EnvVars,
		Timeout:  task.Timeout,
		MetaData: task.MetaData,
		Times: &types.Times{
			ST: time.Now().UnixNano(),
		},
	}

	var stepFnMap = make(map[string]*dag.Vertex)
	for k, step := range task.Steps {
		state, stepFn := task.buildStep(step)
		stepFnMap[step.ID] = dag.NewVertex(step.ID, stepFn)
		task.Steps[k].state = state
	}

	var err error
	defer func() {
		if err != nil {
			logx.Errorln(task.ID, task.workspace, task.scriptDir, err)
		}
	}()

	// 编排步骤: 创建一个有向无环图，图中的每个顶点都是一个作业
	task.graph = dag.New(task.ID)
	for _, step := range task.Steps {
		stepFn, ok := stepFnMap[step.ID]
		if !ok {
			continue
		}
		// 添加顶点以及设置依赖关系
		vertex, err := task.graph.AddVertex(stepFn)
		if err != nil {
			logx.Errorln(err)
			return err
		}
		vertex.WithDeps(task.buildDeps(stepFnMap, step)...)
	}
	// 校验dag图形
	if err = task.graph.Validator(); err != nil {
		logx.Errorln(task.ID, task.workspace, task.scriptDir, err)
		return err
	}

	// 插入数据
	if err = storage.SetTask(task.ID, task.state); err != nil {
		logx.Errorln(task.ID, task.workspace, task.scriptDir, err)
		return err
	}
	queue.PushBack(func() {
		var ctx, cancel = context.WithCancel(context.Background())
		if task.Timeout > 0 {
			ctx, cancel = context.WithTimeoutCause(context.Background(), task.Timeout, exec.ErrTimeOut)
		}
		defer cancel()
		res := task.run(ctx)
		logx.Infoln(task.ID, task.workspace, task.scriptDir, "end of execution")
		if res != nil {
			logx.Infoln(task.ID, task.workspace, task.scriptDir, res)
		}
	})
	return nil
}

func (t *Task) buildStep(step *TaskStep) (*types.TaskStepState, dag.VertexFunc) {
	// 设置缓存中初始状态
	state := t.initStepCache(step)
	// build step
	return state, func(ctx context.Context, taskName, stepName string) error {
		state.State = exec.Running
		state.Times.ST = time.Now().UnixNano()

		if err := storage.SetTaskStep(t.ID, step.ID, state); err != nil {
			logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
			return err
		}

		defer func() {
			if _err := recover(); _err != nil {
				logx.Errorln(t.ID, t.workspace, t.scriptDir, _err)
				state.State = exec.SystemErr
				state.Code = exec.SystemErr
				state.Message = fmt.Sprintf("%v", _err)
				state.Times.ET = time.Now().UnixNano()
				if err := storage.SetTaskStep(t.ID, step.ID, state); err != nil {
					logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
				}
			}
		}()

		// TODO: 执行前
		logx.Infoln(t.ID, t.workspace, t.scriptDir, "started")
		defer func() {
			// TODO: 执行后
			logx.Infoln(t.ID, t.workspace, t.scriptDir, "end")
		}()

		if err := t.execStep(ctx, step, state); err != nil {
			logx.Errorln(err)
			return err
		}
		return nil
	}
}

func (t *Task) buildDeps(stepFnMap map[string]*dag.Vertex, step *TaskStep) []*dag.Vertex {
	var stepFns []*dag.Vertex
	for _, name := range step.DependsOn {
		_stepFn, _ok := stepFnMap[name]
		if !_ok {
			continue
		}
		stepFns = append(stepFns, _stepFn)
	}
	return stepFns
}

func (t *Task) run(ctx context.Context) error {
	defer func() {
		err := recover()
		if err != nil {
			logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		}
	}()

	// 判断当前图形是否挂起
	if t.graph.Paused() {
		t.graph.WaitResume()
	}

	t.state.State = exec.Running
	if err := storage.SetTask(t.ID, t.state); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		return err
	}

	defer func() {
		// 清理工作目录
		t.clear()
		// 结束时间
		t.state.Times.ET = time.Now().UnixNano()
		// 更新数据
		if err := storage.SetTask(t.ID, t.state); err != nil {
			logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		}
	}()

	if err := t.init(); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		t.state.State = exec.SystemErr
		t.state.Message = err.Error()
		return nil
	}

	var state = exec.Stop

	defer func() {
		t.state.State = state
	}()

	if err := t.graph.Run(ctx); err != nil {
		state = exec.SystemErr
		t.state.Message = err.Error()
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
	}
	return nil
}

func (t *Task) init() error {
	if err := utils.EnsureDirExist(t.workspace); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		return err
	}
	if err := utils.EnsureDirExist(t.scriptDir); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		return err
	}
	return nil
}

func (t *Task) clear() {
	if err := os.RemoveAll(t.scriptDir); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
	}
}

func (t *Task) initStepCache(step *TaskStep) *types.TaskStepState {
	var state = &types.TaskStepState{
		ID:             step.ID,
		State:          exec.Pending,
		Message:        "The current step only proceeds if the previous step succeeds.",
		EnvVars:        step.EnvVars,
		Timeout:        step.Timeout,
		DependsOn:      step.DependsOn,
		CommandType:    step.CommandType,
		CommandContent: step.CommandContent,
	}
	if err := storage.SetTaskStep(t.ID, step.ID, state); err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
	}
	return state
}
