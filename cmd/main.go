package main

import (
	"github.com/jilin7105/ebase/cmd/app"
	"log"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
