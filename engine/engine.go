package engine

import (
	"errors"
	"fmt"
	"github.com/Jeffail/tunny"
	"github.com/natessilva/dag"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/osreapi/cache"
	"runtime"
	"sort"
	"sync"
	"time"
)

type ExecJobs struct {
	lock    *sync.Mutex
	results cache.ExecResults
	TaskID  string
	Jobs    []*cache.Jobs
	Val     *cache.Val
	TTL     time.Duration
}

var (
	Pool            *tunny.Pool
	defaultPoolSize = runtime.NumCPU() * 15
)

func NewExecPool() {
	Pool = tunny.NewFunc(defaultPoolSize, worker)
	go func() {
		tt := time.Tick(1 * time.Second)
		for range tt {
			size := Pool.GetSize()
			percent, _ := cpu.Percent(time.Second, false)
			logrus.Infof("pool size %d; cpu usage %.2f", size, percent)
			if percent[0] <= 80 && size <= defaultPoolSize*80 {
				if Pool.QueueLength() >= int64(size) {
					Pool.SetSize(size + 15)
				}
			} else if percent[0] >= 80 && size > defaultPoolSize {
				Pool.SetSize(size - 15)
			}
		}
	}()
}

func Process(taskID string, jobs []*cache.Jobs, ttl time.Duration) {
	// 临时存储
	val := &cache.Val{
		Jobs:  jobs,
		State: cache.Pending,
	}
	// 插入数据
	cache.Set(taskID, val)
	go Pool.Process(ExecJobs{
		lock:   new(sync.Mutex),
		TaskID: taskID,
		Jobs:   jobs,
		Val:    val,
		TTL:    ttl,
	})
}

func Reload() {
	res, err := cache.ListAllData()
	if err != nil {
		logrus.Error(err)
		return
	}
	var count int
	for k, v := range res {
		if v.State == cache.Stopped {
			continue
		}
		count += 1
		var ttl time.Duration
		for _, _v := range v.Jobs {
			ttl += _v.ExecTimeout
		}
		Process(k, v.Jobs, ttl)
	}
	logrus.Infof("%d tasks have been put back on the execution queue", count)
	fmt.Printf("%d tasks have been put back on the execution queue\n", count)
	runtime.GC()
}

func worker(i interface{}) interface{} {
	e, ok := i.(ExecJobs)
	if !ok {
		logrus.Error("input problem")
		return nil
	}
	// 设置过期时间
	e.TTL += 1 * time.Hour
	e.Val.TTL = e.TTL
	e.Val.StartTimes = time.Now().UnixNano()
	e.Val.State = cache.Running
	// 更新数据
	cache.Set(e.TaskID, e.Val)
	// 开始编排jobs
	var d = new(dag.Runner)
	// create a directed graph, where each vertex in the graph
	// is a pipeline job.
	for k, job := range e.Jobs {
		j := job
		step := k
		d.AddVertex(job.Name, func() error {
			return e.execJob(step, j)
		})
	}
	// create the vertex edges from the values configured in the
	// depends_on attribute.
	for _, job := range e.Jobs {
		for _, dep := range job.DependsOn {
			d.AddEdge(dep, job.Name)
		}
	}
	if err := d.Run(); err != nil {
		if err != execErr {
			logrus.Error(err)
			e.results = cache.ExecResults{
				{
					Name:     fmt.Sprintf("%s-init", e.TaskID),
					Step:     -1,
					ExitCode: 255,
					Output:   err.Error(),
				},
			}
		}
	}
	// sort by Step
	sort.Sort(e.results)
	e.Val.Data = e.results
	e.Val.State = cache.Stopped
	e.Val.CompletedTimes = time.Now().UnixNano()
	cache.Set(e.TaskID, e.Val)
	return nil
}

var execErr = errors.New("exec error")

func (e *ExecJobs) execJob(step int, job *cache.Jobs) error {
	_log := logrus.WithFields(logrus.Fields{
		"name": job.Name,
		"task": e.TaskID,
		"step": step,
		"envs": job.EnvVars,
	})
	if job.EnvVars != nil {
		_log.Info("inject external environment variables")
	}
	_log.Info("staring exec job")
	var result = &cache.ExecStatus{
		Name:     job.Name,
		Step:     step,
		ExitCode: 255,
	}
	defer func() {
		e.lock.Lock()
		e.results = append(e.results, result)
		e.lock.Unlock()
	}()
	cmd := &Cmd{
		TaskID:          e.TaskID,
		Step:            step,
		Name:            job.Name,
		Shell:           job.CommandType,
		Content:         job.CommandContent,
		ExternalEnvVars: job.EnvVars,
		ExecTimeout:     job.ExecTimeout,
		Log:             _log,
	}
	if err := cmd.Initial(); err != nil {
		_log.Error(err)
		result.Output = err.Error()
		return execErr
	}
	result.ExitCode, result.Output = cmd.ExecScript(e.TaskID)
	if result.ExitCode != 0 {
		_log.Errorf("exit code is not 0 but %d", result.ExitCode)
		return execErr
	}
	return nil
}
