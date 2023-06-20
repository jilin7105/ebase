package ipLimiter

import (
	"sync"

	"github.com/golang/groupcache/lru"
	"golang.org/x/sync/semaphore"
)

type IPLimiter struct {
	sync.Mutex
	IPMap   *lru.Cache
	ipLimit int
}

func NewIPLimiter(ipLimit, maxEntries int) *IPLimiter {
	return &IPLimiter{
		IPMap:   lru.New(maxEntries),
		ipLimit: ipLimit,
	}
}

func (l *IPLimiter) GetLimiter(ip string) *semaphore.Weighted {
	l.Lock()
	defer l.Unlock()

	limiter, exists := l.IPMap.Get(ip)
	if !exists {
		limiter = semaphore.NewWeighted(int64(l.ipLimit)) // 每个IP限制为同时处理5个并发请求
		l.IPMap.Add(ip, limiter)
	}

	return limiter.(*semaphore.Weighted)
}
