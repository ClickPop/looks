package utils

import (
	"fmt"
	"regexp"
	"strings"
)

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

func ReverseStringSlice(slice []string) []string {
	var revSlice []string
	for i := len(slice) - 1; i >= 0; i-- {
		revSlice = append(revSlice, slice[i])
	}
	return revSlice
}

func OxfordJoin(slice []string) string {
	outStr := ""
	complex := false
	punct := ","

	for _, v := range slice {
		if strings.Contains(v, ",") {
			complex = true
			punct = ";"
		}
	}

	if len(slice) == 1 {
		outStr = slice[0]
	} else if len(slice) == 2 && !complex {
		outStr = strings.Join(slice, " and ")
	} else if len(slice) > 2 || (len(slice) == 2 && complex) {
		outStr = strings.Join(slice[0:(len(slice)-1)], fmt.Sprintf("%s ", punct)) + fmt.Sprintf("%s and ", punct) + slice[len(slice)-1]
	}

	return outStr
}