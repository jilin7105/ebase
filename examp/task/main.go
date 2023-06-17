package main

import (
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"log"
)

func main() {

	ebase.Init()
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
