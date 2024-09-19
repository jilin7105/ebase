package main

import (
	"github.com/jilin7105/ebase/cmd/ebasetools/app"
	"log"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
