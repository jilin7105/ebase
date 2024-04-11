package ipLimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestIPLimiter(t *testing.T) {
	ipLimit := 5
	maxEntries := 1000
	ipLimiter := NewIPLimiter(ipLimit, maxEntries)

	// 使用相同的IP进行并发请求
	ip := "127.0.0.1"
	requests := 10

	// 启动多个并发请求
	done := make(chan bool)
	for i := 0; i < requests; i++ {
		go func() {
			limiter := ipLimiter.GetLimiter(ip)
			if !limiter.TryAcquire(1) {
				fmt.Println("Too many concurrent requests for IP:", ip)
				done <- false
				return
			}
			// 模拟请求处理时间
			time.Sleep(time.Second)
			limiter.Release(1)
			done <- true
		}()
	}

	// 等待所有请求完成
	for i := 0; i < requests; i++ {
		<-done
	}

	// 获取IP的并发限制器
	limiter := ipLimiter.GetLimiter(ip)

	// 再次尝试获取并发限制器
	if !limiter.TryAcquire(1) {
		fmt.Println("Too many concurrent requests for IP:", ip)
	} else {
		fmt.Println("Concurrent requests for IP:", ip, "are allowed")
		limiter.Release(1)
	}
}
