package testutils

import (
	"fmt"
	"log"
)

func Logger(context string) func(string, ...interface{}) {
	return func(format string, v ...interface{}) {
		f := fmt.Sprintf("[%10s]: %s", context, format)
		log.Printf(f, v...)
	}
}
