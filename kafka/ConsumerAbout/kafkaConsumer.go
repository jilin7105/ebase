package ConsumerAbout

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	"github.com/xdg-go/scram"
	"log"
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
	saramaConfig.Consumer.MaxProcessingTime = time.Duration(config.MaxProcessingTime) * time.Millisecond
	saramaConfig.Consumer.Group.Session.Timeout = time.Duration(config.SessionTimeout) * time.Millisecond
	saramaConfig.Consumer.Group.Heartbeat.Interval = time.Duration(config.HeartbeatInterval) * time.Millisecond

	if config.SASL_Enable {

		saramaConfig.Net.SASL.Enable = config.SASL_Enable
		saramaConfig.Net.SASL.User = config.SASL_User
		saramaConfig.Net.SASL.Password = config.SASL_Password
		saramaConfig.Net.SASL.Handshake = config.SASL_Handshake
		saramaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(config.SASL_Mechanism)

		if saramaConfig.Net.SASL.Mechanism == sarama.SASLTypeSCRAMSHA512 {
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
			//saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		} else if saramaConfig.Net.SASL.Mechanism == sarama.SASLTypeSCRAMSHA256 {
			saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }

		}

	}

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
				log.Println("这表示消费者应该停止")
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

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}
func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}
func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
