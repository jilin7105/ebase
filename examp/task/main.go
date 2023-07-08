package main

import (
	"context"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"log"
)

//使用go run main.go  启动测试服务
func main() {

	ebase.Init()
	redisc := ebase.GetRedis("test")
	redisc.Do(context.Background(), "get", "info")
	eb := ebase.GetEbInstance()
	s, err := eb.GetTaskServer()
	//https://github.com/go-co-op/gocron 使用这个库
	if err != nil {
		log.Panicln(err.Error())
	}
	//注册定时任务
	s.Every(1).Minute().Do(func() {
		log.Println(1)
	})
	eb.Run()
}
