package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	//"net"
	"time"
)

func InitRpcService(config config.Config) *grpc.Server {
	s := grpc.NewServer(

		//链式调用
		grpc.ChainUnaryInterceptor(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

				ti := time.Now()
				newebaserequestID := uuid.New().String()
				//此处启动链路追踪
				if config.GrpcServer.TraceTracking {
					oldebaserequestID := getMataid(ctx)
					ctx = context.WithValue(ctx, "eb-grpc-request-id", newebaserequestID)
					logger.Info("grcp Tractracking log : [oldebaserequestID:%s;newebaserequestID:%s]", oldebaserequestID, newebaserequestID)
				}

				resp, err = handler(ctx, req)
				logger.Info("Tractracking log : newebaserequestID:[%s]  接口耗时 %v", newebaserequestID, time.Since(ti))
				return resp, err
			},
		),
	)
	return s
}

//获取头id
func getMataid(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	EbaseRequestId := md.Get("eb-grpc-request-id")
	if len(EbaseRequestId) > 0 {
		return EbaseRequestId[0]
	} else {
		return ""
	}

}

func setMataid(ctx context.Context, id string) context.Context {
	md := metadata.Pairs("eb-grpc-request-id", id)
	// 创建带metadata的上下文
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
