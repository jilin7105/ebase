package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase/logger"
	"log"
)

// 定义消费者组的处理器
type ConsumerGroupHandler struct {
	SetupFunc      func() error
	CleanupFunc    func() error
	ActionMessages func(*sarama.ConsumerMessage) error
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
		log.Printf("从主题 %s 收到消息: %s\n", msg.Topic, string(msg.Value))
		if h.ActionMessages != nil {
			err := h.ActionMessages(msg)
			if err != nil {
				logger.Info("处理消息失败:%s", err.Error())
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

func SetActionMessages(p func(*sarama.ConsumerMessage) error) Option {
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
