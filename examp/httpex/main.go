package main

import (
	"fmt"
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

	r.GET("/ping", func(context *gin.Context) {
		value, exists := context.Get("EbaseRequestID")
		if exists {
			fmt.Println("requestID:", value)
		}
	})

	eb.Run()
}
