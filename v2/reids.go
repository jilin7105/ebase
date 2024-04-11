package ebase

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

func (e *Eb) initRedis() {
	// 初始化所有Redis连接
	for _, redisConfig := range e.Config.Redis {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
			PoolSize: redisConfig.PoolSize,
		})

		// 将Redis连接保存到Eb结构体中
		e.Redis[redisConfig.Name] = rdb
	}
}
