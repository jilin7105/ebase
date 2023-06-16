package ebase

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// GetDB 提供名字获取数据库连接
func (e *Eb) GetDB(name string) (*gorm.DB, error) {
	db, ok := e.DBs[name]
	if !ok {
		return nil, errors.New("no such database: " + name)
	}
	return db, nil
}

// GetRedis 提供名字获取Redis客户端
func (e *Eb) GetRedis(name string) (*redis.Client, error) {
	client, ok := e.Redis[name]
	if !ok {
		return nil, errors.New("no such redis client: " + name)
	}
	return client, nil
}
