package main

import (
	"log"

	"github.com/clickpop/looks-cli/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
