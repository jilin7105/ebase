package main

import (
	"log"

	"github.com/jilin7105/ebase/cmd/ebasetools"
)

func main() {
	if err := ebasetools.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
