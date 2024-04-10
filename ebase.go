package ebase

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/kafka/ConsumerAbout"
	"github.com/jilin7105/ebase/kafka/ProducerAbout"
	"github.com/jilin7105/ebase/logger"
	ebasegrpc "github.com/jilin7105/ebase/server/grpc"
	ebasehttp "github.com/jilin7105/ebase/server/http"
	"github.com/jilin7105/ebase/task"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Eb struct {
	cxt            context.Context
	ConfigFileName string
	Config         config.Config
	DBs            map[string]*gorm.DB
	Redis          map[string]*redis.Client
	kafkaProducer  map[string]*ProducerAbout.KafkaProducer
	serviceTask    *gocron.Scheduler
	serciceHttp    *gin.Engine
	projectPath    string
	kafkaConsumer  map[string]*ConsumerAbout.KafkaConsumer
	grpcServer     *grpc.Server
	regfunc        func() error
	heartbeatPush  func() error
	stopFunc       func()
}

// 定义全局的Eb实例
var ebInstance = &Eb{
	cxt:           context.Background(),
	DBs:           map[string]*gorm.DB{},
	Redis:         map[string]*redis.Client{},
	kafkaProducer: map[string]*ProducerAbout.KafkaProducer{},
}

func SetProjectPath(path string) {
	ebInstance.projectPath = path
}

// 在init函数中初始化全局的Eb实例
func Init() {

	//用于兼容只是用 二进制文件

	if ebInstance.projectPath == "" {
		filePath, err := os.Getwd()
		if err != nil {
			return
		}

		ebInstance.projectPath = filePath
	}

	ebInstance.ParseFlags()
	ebInstance.LoadConfig()
	// 从配置中设置日志级别和日志文件
	logger.SetLogLevel(ebInstance.Config.LogLevel)
	logger.SetLogFile(ebInstance.projectPath + "/" + ebInstance.Config.LogFile)
	ebInstance.initRedis()
	ebInstance.initMysql()
	ebInstance.InitKafkaProducer()
	ebInstance.initServer()

	//开启心跳检测，服务注册
	ebInstance.initmicro()

	//链路追踪
	ebInstance.initLinkTracking()
}

func (e *Eb) initServer() {
	switch e.Config.AppType {
	case "HTTP":
		e.serciceHttp = ebasehttp.InitHttp(e.Config)
		logger.Info("--------------------http服务器初始化------------------")
		// 创建HTTP服务
	case "gRPC":
		e.grpcServer = ebasegrpc.InitRpcService(e.Config)
		// 创建gRPC服务
	case "Task":
		e.serviceTask = task.InitTaskServer()
		logger.Info("--------------------定时任务初始化------------------")
		// 创建任务服务
	case "Kafka":
		var err error
		e.kafkaConsumer, err = ConsumerAbout.NewKafkaConsumers(e.Config)
		if err != nil {
			logger.Error("kafka 初始化失败", err)
		}
	default:
		logger.Error("unknown appType: %v", e.Config.AppType)

	}
}
func (e *Eb) ParseFlags() {
	flag.StringVar(&e.ConfigFileName, "i", "config.yml", "The name of the config file.")
	flag.Parse()
}

func (e *Eb) LoadConfig() {
	data, err := ioutil.ReadFile(e.projectPath + "/" + e.ConfigFileName)
	if err != nil {
		log.Fatalf("failed to read config file %s: %v", e.ConfigFileName, err)
	}

	err = yaml.Unmarshal(data, &e.Config)
	if err != nil {
		log.Fatalf("failed to unmarshal config file %s: %v", e.ConfigFileName, err)
	}
}

func (e *Eb) Run() {
	//创建监听退出chan
	c := make(chan os.Signal, 1)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go e.stop(c) //退出方法
	switch e.Config.AppType {
	case "HTTP":
		logger.Info("--------------------http启动------------------")
		logger.Info("port", e.Config.HttpGin.Port)
		e.serciceHttp.Run(fmt.Sprintf(":%v", e.Config.HttpGin.Port))
		// 创建HTTP服务
	case "gRPC":
		logger.Info("--------------------grpc启动------------------")
		e.grpcRun()
		// 创建gRPC服务
	case "Task":
		logger.Info("--------------------定时任务启动------------------")
		e.serviceTask.StartBlocking()
		// 创建任务服务
	case "Kafka":
		e.kafkaRun()
	default:
		log.Fatalf("unknown appType: %v", e.Config.AppType)
	}

}

func (e *Eb) stop(c chan os.Signal) {

	for s := range c {

		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info("exit signal: %s", s.String())
			if e.stopFunc != nil {
				e.stopFunc()
			}
			os.Exit(0)
		default:
			logger.Info("exit signal: %s", s.String())
		}

	}
}

// 提供一个函数来获取全局的Eb实例
func GetEbInstance() *Eb {
	return ebInstance
}
