package ebasehttp

import (
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestInitHttp(t *testing.T) {
	cfg := config.Config{
		HttpGin: struct {
			Port               int  `yaml:"port"`
			AppendPprof        bool `yaml:"appendPprof"`
			IPConcurrencyLimit int  `yaml:"iPConcurrencyLimit"`
		}{
			AppendPprof: true,
		},
	}

	router := InitHttp(cfg)
	if router == nil {
		t.Error("Expected gin.Engine, got nil")
	}

	// Test pprof routes
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/debug/pprof/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
	}

	// Test middleware
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	//
	//c, _ := router.CreateTestContext(w)
	//if c.MustGet("EbaseRequestID") == nil {
	//	t.Error("Expected context value 'EbaseRequestID', got nil")
	//}
}

func TestHttpServer(t *testing.T) {
	// 创建测试用的配置
	httpConfig := config.Config{
		HttpGin: struct {
			Port               int  `yaml:"port"`
			AppendPprof        bool `yaml:"appendPprof"`
			IPConcurrencyLimit int  `yaml:"iPConcurrencyLimit"`
		}{
			AppendPprof:        true,
			IPConcurrencyLimit: 5,
		},
	}

	// 初始化HTTP服务
	engine := InitHttp(httpConfig)
	engine.GET("/test", func(context *gin.Context) {
		time.Sleep(4 * time.Second)
		context.JSON(http.StatusOK, "ok")
	})
	// 启动HTTP服务
	server := httptest.NewServer(engine)
	defer server.Close()

	// 发起并发请求
	wg := sync.WaitGroup{}
	concurrency := 6
	for i := 1; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 发起GET请求
			resp, err := http.Get(server.URL + "/test")
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			defer resp.Body.Close()

			// 检查响应状态码和内容
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}
	time.Sleep(1 * time.Second)
	// 测试超出限制的请求
	resp, err := http.Get(server.URL + "/test")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	// 等待并发请求完成
	wg.Wait()

}
