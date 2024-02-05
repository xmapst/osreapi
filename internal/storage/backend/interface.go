package backend

import (
	"errors"
	"time"

	"github.com/xmapst/osreapi/internal/storage/types"
)

// bucket 类似表名
// key 类似主键, key可按/拆分

type ITaskCache interface {
	TaskList(table, prefix string) (res types.TaskStates, err error)
	TaskDetail(table, task string) (res *types.TaskState, err error)
	TaskStepList(table, task string) (res types.TaskStepStates, err error)
	TaskStepDetail(table, task, step string) (res *types.TaskStepState, err error)
	TaskStepLogList(table, task, step string) (res types.TaskStepLogs, err error)
	SetTask(table, task string, val *types.TaskState) error
	SetTaskStep(table, task, step string, val *types.TaskStepState) error
	SetTaskStepLog(table, task, step string, line int64, val *types.TaskStepLog) error
	Close() (err error)
	Name() string
}

type Value struct {
	TTL   time.Duration
	Value []byte
}

func SafeCopy(des, src []byte) []byte {
	if len(des) < len(src) {
		des = make([]byte, len(src))
	} else {
		des = des[:len(src)]
	}
	copy(des, src)
	return des
}

var ErrNotExist = errors.New("does not exist")
