package main

import (
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase"
	"log"
)

// 定义消费者组的处理器
type consumerGroupHandler struct {
}

// Setup 在新会话开始之前，ConsumeClaim 之前运行
func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

// Cleanup 在会话结束时运行，所有的 ConsumeClaim goroutines 都已经退出
func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim 必须启动一个消费者循环，处理 ConsumerGroupClaim 的 Messages()
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("从主题 %s 收到消息: %s\n", msg.Topic, string(msg.Value))
		// 在这里处理你的消息
		// 标记消息已处理
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	ebase.Init()
	eb := ebase.GetEbInstance()
	eb.RegisterKafkaHandle("Consumer1", consumerGroupHandler{})
	eb.RegisterKafkaHandle("Consumer2", consumerGroupHandler{})
	eb.Run()
}
