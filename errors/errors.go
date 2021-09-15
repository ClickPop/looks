package errors

import (
	"log"
	"os"
)

func HandleError(err error, callbacks ...func(err error) error) error {
	if err != nil {
		log.Fatalln(err)
		for i := 0; i < len(callbacks); i++ {
			callback := callbacks[i]
			val := callback(err)
			if val != nil {
				return val
			}
		}
	}
	return err
}

func Exit(_ error) error {
	os.Exit(1)
	return nil
}
