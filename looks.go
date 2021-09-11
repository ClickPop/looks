package main

import (
	"fmt"
	"log"
	"os"

	"github.com/clickpop/looks/generator"
	"github.com/clickpop/looks/man"
	"github.com/clickpop/looks/utils"
	"github.com/clickpop/looks/utils/config"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Please specify a command. To see a list of command run looks --help")
	}
	command := args[0];
	configFile, err := utils.ParseArgs(args, "--config-file")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	switch command {
	case "gen":
		fallthrough
	case "generate":
		generator.Generate(&conf, nil)
	case "--help":
		fallthrough
	case "-h":
		fmt.Println(man.ROOT_HELP_MENU)
	}
}
