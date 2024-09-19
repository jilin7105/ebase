package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/kafka/ConsumerAbout"
	"log"
	"os"
)

func Action(msg *sarama.ConsumerMessage, ctx context.Context) error {
	log.Printf("从主题 %s ", ctx.Value(""))
	log.Printf("从主题 %s 收到消息: %s\n", msg.Topic, string(msg.Value))
	return nil
}

func setup() error {
	return nil
}

// 使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()
	eb.EasyRegisterKafkaHandle("ab_test", ConsumerAbout.SetActionMessages(Action), ConsumerAbout.SetSetup(setup))

	eb.Run()
	select {}
}
