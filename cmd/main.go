package main

import (
	"ebasetools/app"
	"log"
)

func main() {
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
