package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/util/EbaseGinResponse"
	"log"
	"os"
)

// 使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()

	r, err := eb.GetHttpServer()
	//gin 库
	if err != nil {
		log.Panicln(err.Error())
	}
	eb.SetStopFunc(func() {
		//time.Sleep(15 * time.Second)
		log.Println("stop")
	})
	r.GET("/ping", func(context *gin.Context) {
		EbaseGinResponse.Error(context, 201, fmt.Errorf("info %s err", "test"), "")
	})

	eb.Run()
}
