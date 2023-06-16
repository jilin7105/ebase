package ebase

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	AppType string `yaml:"appType"`
}

type Eb struct {
	ConfigFileName string
	Config         Config
}

// 定义全局的Eb实例
var ebInstance *Eb

// 在init函数中初始化全局的Eb实例
func init() {
	ebInstance = &Eb{}
	ebInstance.ParseFlags()
	ebInstance.LoadConfig()
}

func (e *Eb) ParseFlags() {
	flag.StringVar(&e.ConfigFileName, "i", "config.yaml", "The name of the config file.")
	flag.Parse()
}

func (e *Eb) LoadConfig() {
	data, err := ioutil.ReadFile(e.ConfigFileName)
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
		// 创建HTTP服务
	case "gRPC":
		// 创建gRPC服务
	case "Task":
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
