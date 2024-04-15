package ebasehttp

import (
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/helpfunc"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase/server/ipLimiter"
	"github.com/jilin7105/ebase/util/LinkTracking"
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
		requestID := helpfunc.CreateRequestId()
		if c.GetHeader("EbaseRequestID") == "" {

			c.Header("EbaseRequestID", requestID)
		} else {
			requestID = c.GetHeader("EbaseRequestID")
		}

		c.Set("EbaseRequestID", requestID)
		Span := c.Request.URL
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		elapsedTime := time.Since(startTime)

		data, err := LinkTracking.NewLinkTrackLogData(
			LinkTracking.LinkTrackID(requestID),
			LinkTracking.LinkTrackTime(time.Now().Format("2006-01-02 15:04:05")),
			LinkTracking.LinkTrackActionTime(elapsedTime.String()),
			LinkTracking.LinkTrackSpan(Span.String()),
		)
		if err != nil {
			logger.Error("LinkTracking error: %v", err)
		}
		data.Send()
		// Log request details
		logger.Info("Request ID: %s, Path: %s, Time: %v", requestID, c.Request.URL.Path, elapsedTime)
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
