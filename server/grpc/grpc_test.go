package grpc

import (
	"context"
	"github.com/jilin7105/ebase/config"
	"google.golang.org/grpc/metadata"
	"log"
	"testing"
)

func TestInitRpcService(t *testing.T) {
	cfg := config.Config{
		GrpcServer: struct {
			Port          int  `yaml:"port"`
			TraceTracking bool `yaml:"traceTracking"`
		}{
			Port:          8080,
			TraceTracking: true,
		},
	}

	server := InitRpcService(cfg)
	if server == nil {
		t.Error("Expected grpc server, got nil")
	}

	// Test interceptor
	interceptor := server.GetServiceInfo()
	log.Println(interceptor)
	if interceptor == nil {
		t.Error("Expected unary interceptor, got nil")
	}

	// Test metadata functions
	ctx := context.Background()
	ctx = setMataid(ctx, "test-id")
	id := getMataid(ctx)
	if id != "test-id" {
		t.Errorf("Expected 'test-id', got %v", id)
	}

	md, _ := metadata.FromOutgoingContext(ctx)
	if len(md.Get("eb-grpc-request-id")) == 0 {
		t.Error("Expected metadata with 'eb-grpc-request-id', got none")
	}
}
