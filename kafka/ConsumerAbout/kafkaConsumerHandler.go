package ConsumerAbout

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase/util/LinkTracking"
	"time"
)

// 定义消费者组的处理器
type ConsumerGroupHandler struct {
	SetupFunc      func() error
	CleanupFunc    func() error
	ActionMessages func(*sarama.ConsumerMessage, context.Context) error
}

// Setup 在新会话开始之前，ConsumeClaim 之前运行
func (h ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	if h.SetupFunc != nil {
		return h.SetupFunc()
	}
	return nil
}

// Cleanup 在会话结束时运行，所有的 ConsumeClaim goroutines 都已经退出
func (h ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	if h.CleanupFunc != nil {
		return h.CleanupFunc()
	}
	return nil
}

// ConsumeClaim 必须启动一个消费者循环，处理 ConsumerGroupClaim 的 Messages()
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		requestID := ""

		if h.ActionMessages != nil {

			//从 kafka 数据header 中获取requestID
			headrs := msg.Headers
			for _, RecordHeader := range headrs {
				if string(RecordHeader.Key) == "EbaseRequestID" {
					requestID = string(RecordHeader.Value)
				}

			}
			ctx := context.WithValue(context.Background(), "EbaseRequestID", requestID)

			// Start timer
			startTime := time.Now()

			err := h.ActionMessages(msg, ctx)
			if err != nil {
				logger.Info("处理消息失败:%s", err.Error())
			}

			if LinkTracking.GetIsOpen() && requestID != "" {
				// Calculate request time
				elapsedTime := time.Since(startTime)

				data, err := LinkTracking.NewLinkTrackLogData(
					LinkTracking.LinkTrackID(requestID),
					LinkTracking.LinkTrackTime(time.Now().Format("2006-01-02 15:04:05")),
					LinkTracking.LinkTrackActionTime(elapsedTime.String()),
					LinkTracking.LinkTrackSpan(msg.Topic),
				)
				if err != nil {
					logger.Error("发送链路追踪数据失败:%s", err.Error())
				}
				data.Send()
			}
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

type Option func(*ConsumerGroupHandler)

func SetSetup(p func() error) Option {
	return func(s *ConsumerGroupHandler) {
		s.SetupFunc = p
	}
}
func SetCleanupFunc(p func() error) Option {
	return func(s *ConsumerGroupHandler) {
		s.CleanupFunc = p
	}
}

func SetActionMessages(p func(*sarama.ConsumerMessage, context.Context) error) Option {
	return func(s *ConsumerGroupHandler) {
		s.ActionMessages = p
	}
}

func NewConsumerHandler(options ...func(*ConsumerGroupHandler)) (*ConsumerGroupHandler, error) {
	handler := ConsumerGroupHandler{}
	for _, option := range options {
		option(&handler)
	}
	if handler.ActionMessages == nil {
		return nil, fmt.Errorf("未设置处理消息函数,创建失败")
	}
	return &handler, nil
}
