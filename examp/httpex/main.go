package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"log"
	"os"
	"path/filepath"
	"time"
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

//使用go run main.go  启动测试服务
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
		time.Sleep(1 * 5 * time.Second)
		log.Println("stop")
	})
	r.GET("/ping", func(context *gin.Context) {
		value, exists := context.Get("EbaseRequestID")
		if exists {
			fmt.Println("requestID:", value)
		}
	})

	eb.Run()
}
