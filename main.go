package main

import (
	"log"

	"github.com/clickpop/looks/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
