package ebasehttp

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"net/http/pprof"
	"time"
)

func InitHttp(config config.Config) *gin.Engine {
	r := gin.Default()
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
