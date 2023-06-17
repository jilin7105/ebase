package ebase

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/jilin7105/ebase/config"
	"github.com/jilin7105/ebase/logger"
	ebasehttp "github.com/jilin7105/ebase/server/http"
	"github.com/jilin7105/ebase/task"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type Eb struct {
	ConfigFileName string
	Config         config.Config
	DBs            map[string]*gorm.DB
	Redis          map[string]*redis.Client
	serviceTask    *gocron.Scheduler
	serciceHttp    *gin.Engine
	projectPath    string
}

// 定义全局的Eb实例
var ebInstance *Eb

// 在init函数中初始化全局的Eb实例
func Init() {

	ebInstance = &Eb{
		DBs:   map[string]*gorm.DB{},
		Redis: map[string]*redis.Client{},
	}
	_, filePath, _, _ := runtime.Caller(1)
	log.Println(filePath)
	ebInstance.projectPath = filepath.Dir(filePath)
	log.Println(ebInstance.projectPath)
	ebInstance.ParseFlags()
	ebInstance.LoadConfig()
	// 从配置中设置日志级别和日志文件
	logger.SetLogLevel(ebInstance.Config.LogLevel)
	logger.SetLogFile(ebInstance.projectPath + "/" + ebInstance.Config.LogFile)
	ebInstance.initRedis()
	ebInstance.initMysql()
	ebInstance.initServer()
}

func (e *Eb) initServer() {
	switch e.Config.AppType {
	case "HTTP":
		e.serciceHttp = ebasehttp.InitHttp(e.Config)
		// 创建HTTP服务
	case "gRPC":
		// 创建gRPC服务
	case "Task":
		e.serviceTask = task.InitTaskServer()
		logger.Info("--------------------定时任务启动------------------")
		// 创建任务服务
	case "Kafka":
		// 创建Kafka服务
	default:
		log.Fatalf("unknown appType: %v", e.Config.AppType)
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
	switch e.Config.AppType {
	case "HTTP":
		e.serciceHttp.Run(fmt.Sprintf(":%d", e.Config.HttpGin.Port))
		// 创建HTTP服务
	case "gRPC":
		// 创建gRPC服务
	case "Task":
		e.serviceTask.StartBlocking()
		logger.Info("--------------------定时任务启动------------------")
		// 创建任务服务
	case "Kafka":
		// 创建Kafka服务
	default:
		log.Fatalf("unknown appType: %v", e.Config.AppType)
	}
}

// 提供一个函数来获取全局的Eb实例
func GetEbInstance() *Eb {

	return ebInstance
}
