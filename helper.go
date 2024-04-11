package ebase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/jilin7105/ebase/kafka/ConsumerAbout"
	"github.com/jilin7105/ebase/kafka/ProducerAbout"
	"github.com/jilin7105/ebase/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net"
)

// GetDB 提供名字获取数据库连接
func GetDB(name string) *gorm.DB {
	db, ok := ebInstance.DBs[name]
	if !ok {
		return nil
	}
	return db
}

// GetRedis 提供名字获取Redis客户端
func GetRedis(name string) *redis.Client {
	client, ok := ebInstance.Redis[name]
	if !ok {
		return nil
	}
	return client
}

// GetKafka 获取Kafka生产者
func GetKafka(name string) *ProducerAbout.KafkaProducer {
	client, ok := ebInstance.kafkaProducer[name]
	if !ok {
		return nil
	}
	return client
}

// GetEs 获取ES客户端
func GetEs(name string) *elasticsearch.Client {
	client, ok := ebInstance.ES[name]
	if !ok {
		return nil
	}
	return client
}

// GetMongo 获取Mongo客户端
func GetMongo(name string) *mongo.Client {
	client, ok := ebInstance.Mongo[name]
	if !ok {
		return nil
	}
	return client
}

func (e *Eb) GetTaskServer() (*gocron.Scheduler, error) {
	service := e.serviceTask

	if service == nil {
		return nil, errors.New("未初始化定时任务服务(Task)，请检测服务类型 ")
	}

	return service, nil
}

func (e *Eb) GetGrpcServer() (*grpc.Server, error) {
	service := e.grpcServer

	if service == nil {
		return nil, errors.New("未初始化定时任务服务(Grpc)，请检测服务类型 ")
	}

	return service, nil
}

func (e *Eb) GetHttpServer() (*gin.Engine, error) {
	service := e.serciceHttp

	if service == nil {
		return nil, errors.New("未初始化定时任务服务(http)，请检测服务类型 ")
	}

	return service, nil
}

// 注册kafka 消费
func (e *Eb) RegisterKafkaHandle(name string, handle sarama.ConsumerGroupHandler) {
	if _, ok := e.kafkaConsumer[name]; ok {
		e.kafkaConsumer[name].RegisterHandle(handle)
	}
}

// 启动kafka 消费
func (e *Eb) kafkaRun() {
	for _, consumer := range e.kafkaConsumer {
		go func(consumer *ConsumerAbout.KafkaConsumer, ctx context.Context) {
			log.Printf("Starting Kafka consumer: %s", consumer.Name)
			if err := consumer.Consume(ctx); err != nil {
				log.Printf("Error consuming from Kafka (consumer: %s): %v", consumer.Name, err)
			}
		}(consumer, e.cxt)
	}
}

// 启动kafka 消费
func (e *Eb) grpcRun() {
	port := e.Config.GrpcServer.Port
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	e.grpcServer.Serve(ln)
}

// 自动加载配置文件
func (e *Eb) SelfLoadConfig(out interface{}) error {
	data, err := ioutil.ReadFile(e.projectPath + "/" + e.ConfigFileName)
	if err != nil {
		log.Fatalf("failed to read config file %s: %v", e.ConfigFileName, err)
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		log.Fatalf("failed to unmarshal config file %s: %v", e.ConfigFileName, err)
		return err
	}
	return err
}

// 写入退出回调信息
func (e *Eb) SetStopFunc(f func()) {
	e.stopFunc = f
}

func (e *Eb) EasyRegisterKafkaHandle(name string, options ...func(*ConsumerAbout.ConsumerGroupHandler)) {
	if _, ok := e.kafkaConsumer[name]; ok {
		handle, err := ConsumerAbout.NewConsumerHandler(options...)
		if err != nil {
			logger.Info("kafka handle init error name[%s]  err:%s", name, err.Error())
			return
		}
		e.kafkaConsumer[name].RegisterHandle(handle)
	}
}
