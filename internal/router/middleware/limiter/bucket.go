package limiter

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens   int64
	lock     sync.Mutex
	rate     int64 // 每秒加入的令牌数
	lastTime int64
}

func newBucket(rate int64) *Bucket {
	b := &Bucket{
		rate: rate,
	}

	return b
}

// IsAccept 是否接受请求
func (b *Bucket) IsAccept() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	now := time.Now().Unix()
	b.tokens = b.tokens + (now-b.lastTime)*b.rate

	if b.tokens >= 1 {
		b.tokens = 1
	}
	b.lastTime = now

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}
