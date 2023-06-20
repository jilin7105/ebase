package ebasehttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/server/ipLimiter"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

var ipConcurrencyLimit int
var limiter *ipLimiter.IPLimiter

func InitHttp(config config.Config) *gin.Engine {
	r := gin.Default()
	//如果 ip限制大于0 默认使用ip 并发限制中间键
	ipConcurrencyLimit = config.HttpGin.IPConcurrencyLimit
	if ipConcurrencyLimit > 0 {
		limiter = ipLimiter.NewIPLimiter(ipConcurrencyLimit, 1000)
		r.Use(limitByIP())
	}

	//增加接口请求日志，增加requestID
	r.Use(func(c *gin.Context) {
		// Generate a unique ID for this request
		requestID := uuid.New().String()

		// Add the request ID to the context
		c.Set("EbaseRequestID", requestID)

		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate request time
		elapsedTime := time.Since(startTime)
		fmt.Sprintf("Request ID: %s, Path: %s, Time: %v", requestID, c.Request.URL.Path, elapsedTime)
		// Log request details
		//logger.Info("Request ID: %s, Path: %s, Time: %v", requestID, c.Request.URL.Path, elapsedTime)
	})

	if config.HttpGin.AppendPprof {
		appendPprof(r)
	}

	//// Your routes and middleware go here
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	return r
}

func appendPprof(r *gin.Engine) {
	r.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	r.GET("/debug/pprof/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
	r.GET("/debug/pprof/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
	r.GET("/debug/pprof/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
	r.GET("/debug/pprof/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	r.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	r.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	r.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	r.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))
}

func limitByIP() gin.HandlerFunc {

	return LimitByIP(limiter)
}

/**
ip并发控制
*/
func LimitByIP(limiter *ipLimiter.IPLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ipLimiter := limiter.GetLimiter(ip)
		if !ipLimiter.TryAcquire(1) {

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}

		defer ipLimiter.Release(1)

		c.Next()
	}
}
