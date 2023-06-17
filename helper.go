package ebase

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// GetDB 提供名字获取数据库连接
func GetDB(name string) (*gorm.DB, error) {
	db, ok := ebInstance.DBs[name]
	if !ok {
		return nil, errors.New("no such database: " + name)
	}
	return db, nil
}

// GetRedis 提供名字获取Redis客户端
func GetRedis(name string) (*redis.Client, error) {
	client, ok := ebInstance.Redis[name]
	if !ok {
		return nil, errors.New("no such redis client: " + name)
	}
	return client, nil
}

func (e *Eb) GetTaskServer() (*gocron.Scheduler, error) {
	service := e.serviceTask

	if service == nil {
		return nil, errors.New("未初始化定时任务服务(Task)，请检测服务类型 ")
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
