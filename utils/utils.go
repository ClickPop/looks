package utils

import (
	"errors"
	"regexp"
	"strings"
)

func ParseArgs(args []string, argVal string) (string, error) {
	for i, v := range args {
		if v == argVal {
			if i == len(args)-1 {
				return "nil", errors.New("No value given for flag")
			}
			return args[i+1], nil
		}
	}
	return "", nil
}

func Contains(slice []string, haystack string) bool {
	rVal := false

	for _, v := range slice {
		if v == haystack {
			rVal = true
		}
	}

	return rVal
}

func TransformName(name string) string {
	rx := regexp.MustCompile("(_|-)")
	newName := strings.Title(rx.ReplaceAllString(name, " "))
	return newName
}