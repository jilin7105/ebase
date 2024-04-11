package ebase

import (
	"github.com/Shopify/sarama"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/kafka/ProducerAbout"
	"github.com/jilin7105/ebase/logger"
	"time"
)

func newKafkaProducer(config config.KafkaProducerConfig) (*ProducerAbout.KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	// 设置默认值
	if config.Timeout == 0 {
		config.Timeout = 5000
	}
	saramaConfig.Net.DialTimeout = time.Duration(config.Timeout) * time.Millisecond

	if config.BatchSize == 0 {
		config.BatchSize = 500
	}
	saramaConfig.Producer.Flush.Messages = config.BatchSize

	if config.BatchTime == 0 {
		config.BatchTime = 5000
	}
	saramaConfig.Producer.Flush.Frequency = time.Duration(config.BatchTime) * time.Millisecond

	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	saramaConfig.Producer.Retry.Max = config.MaxRetries

	if config.RetryBackoff == 0 {
		config.RetryBackoff = 100
	}
	saramaConfig.Producer.Retry.Backoff = time.Duration(config.RetryBackoff) * time.Millisecond

	// 其他设置
	saramaConfig.Producer.Return.Successes = config.ReturnSuccesses

	if config.WaitForAll {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	} else {
		saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	}
	//"gzip"、、"lz4"、"zstd"、"none"
	switch config.Compression {
	case "gzip":
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		saramaConfig.Producer.Compression = sarama.CompressionSnappy
	case "zstd":
		saramaConfig.Producer.Compression = sarama.CompressionZSTD
	case "none":
		saramaConfig.Producer.Compression = sarama.CompressionNone
	default:
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	}

	if config.NewManualPartitioner {
		saramaConfig.Producer.Partitioner = sarama.NewManualPartitioner
	}
	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	var kp = ProducerAbout.KafkaProducer{
		Sp:    &producer,
		Topic: config.Topic,
	}
	return &kp, err

}

func (e *Eb) InitKafkaProducer() {

	for _, config := range e.Config.KafkaProducers {
		Producer, err := newKafkaProducer(config)
		if err != nil {
			logger.Error("kafka 消费者 初始化失败 %s %s", config.Name, err.Error())
		}
		e.kafkaProducer[config.Name] = Producer
	}

}
