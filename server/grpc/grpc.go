package grpc

import (
	"context"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/helpfunc"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase/util/LinkTracking"
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
				// Start timer
				startTime := time.Now()
				requestID := getMataid(ctx)
				if requestID == "" {
					requestID = helpfunc.CreateRequestId()
				}
				ctx = context.WithValue(ctx, "EbaseRequestID", requestID)

				Span := info.FullMethod

				resp, err = handler(ctx, req)

				if LinkTracking.GetIsOpen() && requestID != "" {
					elapsedTime := time.Since(startTime)

					data, err := LinkTracking.NewLinkTrackLogData(
						LinkTracking.LinkTrackID(requestID),
						LinkTracking.LinkTrackTime(time.Now().Format("2006-01-02 15:04:05")),
						LinkTracking.LinkTrackActionTime(elapsedTime.String()),
						LinkTracking.LinkTrackSpan(Span),
					)
					if err != nil {
						logger.Error("LinkTracking error: %v", err)
					}
					data.Send()
				}

				return resp, err
			},
		),
	)
	return s
}

// 获取头id
func getMataid(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	EbaseRequestId := md.Get("EbaseRequestID")
	if len(EbaseRequestId) > 0 {
		return EbaseRequestId[0]
	} else {
		return ""
	}

}

func setMataid(ctx context.Context, id string) context.Context {
	md := metadata.Pairs("EbaseRequestID", id)
	// 创建带metadata的上下文
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
