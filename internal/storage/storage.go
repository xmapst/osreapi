package storage

import (
	"github.com/xmapst/osreapi/internal/storage/backend"
	"github.com/xmapst/osreapi/internal/storage/backend/bolt"
	"github.com/xmapst/osreapi/internal/storage/types"
)

var db backend.ITaskCache

func New(t, d string) (err error) {
	switch t {
	default:
		// bolt
		db, err = bolt.New(d)
	}
	if err != nil {
		return err
	}
	return nil
}

func TaskList() (res types.TaskStates, err error) {
	return db.TaskList(types.TableTask, "")
}

func TaskDetail(task string) (res *types.TaskState, err error) {
	return db.TaskDetail(types.TableTask, task)
}

func TaskStepList(task string) (res types.TaskStepStates, err error) {
	return db.TaskStepList(types.TableStep, task)
}

func TaskStepDetail(task, step string) (res *types.TaskStepState, err error) {
	return db.TaskStepDetail(types.TableStep, task, step)
}

func TaskStepLogList(task, step string) (res types.TaskStepLogs, err error) {
	return db.TaskStepLogList(types.TableLog, task, step)
}

func SetTask(task string, val *types.TaskState) error {
	return db.SetTask(types.TableTask, task, val)
}

func SetTaskStep(task, step string, val *types.TaskStepState) error {
	return db.SetTaskStep(types.TableStep, task, step, val)
}

func SetTaskStepLog(task, step string, line int64, val *types.TaskStepLog) error {
	return db.SetTaskStepLog(types.TableLog, task, step, line, val)
}
func Close() error {
	return db.Close()
}
