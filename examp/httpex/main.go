package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jilin7105/ebase"
	_ "github.com/jilin7105/ebase"
	"log"
)

func main() {

	ebase.Init()
	eb := ebase.GetEbInstance()
	log.Println(eb.Config)
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
