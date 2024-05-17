package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/helpfunc/EBHttpRequest"
	"github.com/jilin7105/ebase/logger"
	"github.com/jilin7105/ebase/util/EbaseGinResponse"
	"github.com/levigross/grequests"
	"log"
	"os"
	"path/filepath"
)

func getExecutableDir() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	executableDir := filepath.Dir(executablePath)
	return executableDir, nil
}

type Myconfig struct {
	//配置文件
	Auto string `json:"auto"`
}

// 使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()
	mycfg := Myconfig{}
	eb.SelfLoadConfig(&mycfg)
	log.Printf("++%v", mycfg)
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
		EBHttpRequest.Get(context, "http://127.0.0.1:9999/f", nil)
		//EbaseGinResponse.OK(context, "pong", "")
		//EbaseGinResponse.PageOK(context, []int{1, 23, 4, 5}, 10, 1, 1, "")
		EbaseGinResponse.Error(context, 201, fmt.Errorf("info %s err", "test"), "")
	})
	r.GET("/f", func(context *gin.Context) {
		EBHttpRequest.Post(context, "http://127.0.0.1:9999/p", &grequests.RequestOptions{
			JSON: map[string]string{
				"name": "jilin",
			},
		})
		EbaseGinResponse.Error(context, 201, fmt.Errorf("info %s err", "test"), "")
	})

	r.POST("/p", func(context *gin.Context) {
		//获取并打印 json 数据
		var a = map[string]string{}
		err := context.BindJSON(&a)
		if err != nil {
			logger.Info("%s", err.Error())
			return
		}

		logger.Info("%++v", a)
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	eb.Run()
}
