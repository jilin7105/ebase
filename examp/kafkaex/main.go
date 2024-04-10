package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/kafka/ConsumerAbout"
	"log"
	"os"
)

//
//// 定义消费者组的处理器
//type consumerGroupHandler struct {
//}
//
//// Setup 在新会话开始之前，ConsumeClaim 之前运行
//func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }
//
//// Cleanup 在会话结束时运行，所有的 ConsumeClaim goroutines 都已经退出
//func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
//
//// ConsumeClaim 必须启动一个消费者循环，处理 ConsumerGroupClaim 的 Messages()
//func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//	for msg := range claim.Messages() {
//		log.Printf("从主题 %s 收到消息: %s\n", msg.Topic, string(msg.Value))
//
//		sess.MarkMessage(msg, "")
//	}
//	return nil
//}

func Action(msg *sarama.ConsumerMessage, ctx context.Context) error {
	log.Printf("从主题 %s ", ctx.Value(""))
	log.Printf("从主题 %s 收到消息: %s\n", msg.Topic, string(msg.Value))
	return nil
}

func setup() error {
	return nil
}

//使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()
	//eb.RegisterKafkaHandle("ab_test", consumerGroupHandler{})

	eb.EasyRegisterKafkaHandle("ab_test", ConsumerAbout.SetActionMessages(Action), ConsumerAbout.SetSetup(setup))

	eb.Run()
	select {}
}
