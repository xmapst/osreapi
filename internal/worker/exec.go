package worker

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/xmapst/osreapi/internal/exec"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/storage"
	"github.com/xmapst/osreapi/internal/storage/types"
)

const ConsoleStart = "OSREAPI::CONSOLE::START"
const ConsoleDone = "OSREAPI::CONSOLE::DONE"

func (t *Task) execStep(ctx context.Context, step *TaskStep, state *types.TaskStepState) error {
	defer func() {
		if _err := recover(); _err != nil {
			logx.Errorln(_err)
			state.Code = exec.SystemErr
			state.Message = fmt.Sprintf("%v", _err)
		}
		state.Times.ET = time.Now().UnixNano()
		switch state.Code {
		case exec.Killed, exec.Timeout, exec.SystemErr:
			state.State = state.Code
		default:
			state.State = exec.Stop
		}
		if _err := storage.SetTaskStep(t.ID, step.ID, state); _err != nil {
			logx.Errorln(_err)
		}
	}()
	var logCh = make(chan string, 1024)
	var cmd, err = exec.New(
		exec.WithLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))),
		exec.WithEnv(append(append(t.EnvVars, step.EnvVars...),
			fmt.Sprintf("TASK_ID=%s", t.ID),
			fmt.Sprintf("TASK_STEP_ID=%s", step.ID)),
		),
		exec.WithShell(step.CommandType),
		exec.WithScript(step.CommandContent),
		exec.WithWorkspace(t.workspace),
		exec.WithScriptDir(t.scriptDir),
		exec.WithTimeout(step.Timeout),
		exec.WithConsoleCh(logCh),
	)
	if err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, err)
		state.Message = err.Error()
		state.Code = exec.SystemErr
		return err
	}
	go t.writeLog(step.ID, logCh)
	state.Code, err = cmd.Run(ctx)
	if err != nil {
		logx.Errorln(t.ID, t.workspace, t.scriptDir, "exit code", state.Code, err)
		state.Message = err.Error()
		return err
	}
	state.Message = "execution succeed"
	return nil
}

func (t *Task) writeLog(stepID string, logCh chan string) {
	var num int64
	// start
	if err := storage.SetTaskStepLog(t.ID, stepID, num, &types.TaskStepLog{
		Timestamp: time.Now().UnixNano(),
		Line:      num,
		Content:   ConsoleStart,
	}); err != nil {
		logx.Warnln(err)
	}
	defer func() {
		// end
		num += 1
		if err := storage.SetTaskStepLog(t.ID, stepID, num, &types.TaskStepLog{
			Timestamp: time.Now().UnixNano(),
			Line:      num,
			Content:   ConsoleDone,
		}); err != nil {
			logx.Warnln(err)
		}
	}()
	// content
	for log := range logCh {
		// TODO: 从输出中获取内容设置到环境变量中心

		num += 1
		if err := storage.SetTaskStepLog(t.ID, stepID, num, &types.TaskStepLog{
			Timestamp: time.Now().UnixNano(),
			Line:      num,
			Content:   log,
		}); err != nil {
			logx.Warnln(err)
		}
	}
}
