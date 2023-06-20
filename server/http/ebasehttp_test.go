package ebasehttp

import (
	"github.com/jilin7105/ebase/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitHttp(t *testing.T) {
	cfg := config.Config{
		HttpGin: struct {
			Port        int  "yaml:\"port\""
			AppendPprof bool "yaml:\"appendPprof\""
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
