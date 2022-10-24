package config

import (
	"log"
	"os"
	"regexp"
	"strings"

	conf "github.com/clickpop/looks/pkg/config"
)

func ValidateConfig(config *conf.Config) {
	validateInput(config.Input.Local.Pathname, config.Input.Local.Filename)
	validateOutput(config.Output.Local.Directory)
}

func validateOutput(path string) {
	if _, err := os.ReadDir(path); err != nil {
		log.Fatalf("issue with output path: %q", err)
	}
}

func validateInput(path string, filename string) {
	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("issue with path to pieces: %q", err)
	}

	if len(dir) < 1 {
		log.Fatalf("directory %s is empty", path)
	} else {
		for _, entry := range dir {
			filePath := path + "/" + entry.Name()
			if entry.IsDir() {
				validateInput(filePath, filename)
			} else {
				checkFileName(entry.Name(), filename)
			}
		}
	}
}

func checkFileName(path string, filename string) {
	expr := strings.ReplaceAll(strings.ReplaceAll(filename, ".", "\\."), "%s", ".*")
	if !regexp.MustCompile(expr).MatchString(path) {
		log.Fatalf("path %s doesn't match filename %s", path, filename)
	}
}
