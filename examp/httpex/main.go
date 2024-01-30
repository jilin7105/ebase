package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
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
		//time.Sleep(15 * time.Second)
		log.Println("stop")
	})
	log.Println(eb.AutoDBs)
	r.GET("/ping", func(context *gin.Context) {

		db := ebase.GetAutoDB("test12")
		if db == nil {
			log.Panicln("1111")
			return
		}
		type User struct {
			UserId int
		}
		var users = []User{}
		db.Table("").Select("user_id").Where("dt = ?", "20240129").First(&users)
		log.Println(users)
	})

	eb.Run()
}
