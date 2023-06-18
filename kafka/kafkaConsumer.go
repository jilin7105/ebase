package kafka

import (
	"context"
	"fmt"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConsumer struct {
	Name     string
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
}

func NewKafkaConsumer(config *config.KafkaConsumerConfig) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	if config.AutoOffsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	saramaConfig.Consumer.MaxWaitTime = time.Duration(config.MaxWaitTime) * time.Millisecond
	saramaConfig.Consumer.Group.Session.Timeout = time.Duration(config.SessionTimeout) * time.Millisecond
	saramaConfig.Consumer.Group.Heartbeat.Interval = time.Duration(config.HeartbeatInterval) * time.Millisecond

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Name:     config.Name,
		consumer: consumer,
		topics:   config.Topics,
	}, nil
}

func NewKafkaConsumers(configs config.Config) (map[string]*KafkaConsumer, error) {
	consumers := map[string]*KafkaConsumer{}
	for _, config := range configs.KafkaConsumers {
		consumer, err := NewKafkaConsumer(config)
		if err != nil {
			return nil, err
		}
		consumers[config.Name] = consumer
	}
	return consumers, nil
}

// 消费消息
func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	// 定义一个消费者组处理器
	//handler := consumerGroupHandler{kc}
	if kc.handler == nil {
		return fmt.Errorf("kafka 未注册 %s ", kc.Name)
	}
	// 使用消费者组从主题中消费消息
	go func() {
		for {
			// `Consume` 应该在一个无限循环中被调用，当服务器端的重新平衡发生时，
			// 消费者会话将需要被重新创建以获取新的声明
			if err := kc.consumer.Consume(ctx, kc.topics, kc.handler); err != nil {
				logger.Info("从Kafka消费时出错 (消费者: %s): %v", kc.Name, err)
			}
			// 检查是否取消了上下文，这表示消费者应该停止
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// 等待上下文取消。这可能是从另一个函数发出的信号，表示这个消费者应该停止。
	<-ctx.Done()

	return kc.consumer.Close()
}

// 关闭消费者
func (kc *KafkaConsumer) Close() error {
	if err := kc.consumer.Close(); err != nil {
		logger.Error("关闭消费者", kc.Name, err)
		return err
	}
	return nil
}

// 注册
func (kc *KafkaConsumer) RegisterHandle(handle sarama.ConsumerGroupHandler) {
	kc.handler = handle
}
