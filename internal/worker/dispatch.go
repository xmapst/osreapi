package worker

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/xmapst/osreapi/internal/deque"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/tunny"
)

var (
	// DefaultSize 默认worker数为cpu核心数的两倍
	DefaultSize = runtime.NumCPU() * 2
	pool        = tunny.NewCallback(DefaultSize)
	queue       = deque.New[func()]()
	taskTotal   int64
)

func init() {
	go dispatch()
}

func SetSize(n int) {
	pool.SetSize(n)
}

func GetSize() int {
	return pool.GetSize()
}

func GetTotal() int64 {
	return atomic.LoadInt64(&taskTotal)
}

func Running() int64 {
	if pool.QueueLength() > int64(pool.GetSize()) {
		return pool.QueueLength() - 1
	}
	return pool.QueueLength()
}

func Waiting() int64 {
	if pool.QueueLength() > int64(pool.GetSize()) {
		return int64(queue.Len()) + 1
	}
	return int64(queue.Len())
}

func StopWait() {
	logx.Info("Waiting for all tasks to complete")
	for queue.Len() != 0 || pool.QueueLength() != 0 {
		time.Sleep(500 * time.Millisecond)
	}
	logx.Info("All tasks completed, normal end")
}

func dispatch() {
	for {
		fn, err := queue.PopFront()
		if err != nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		pool.Submit(fn)
	}
}
