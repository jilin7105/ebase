package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/logger"
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

		esdb := ebase.GetEs("es1")
		if esdb == nil {
			logger.Info("esdb is nil")
		} else {
			info, err := esdb.Info()
			if err != nil {
				logger.Info("esdb.Info() err:%v", err)
			}

			logger.Info("esdb.Info() :%v", info)
		}

		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	eb.Run()
}
