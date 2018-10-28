// +build debug

package delayed

import (
	"fmt"
	"log"
)

func logger(context string) func(string, ...interface{}) {
	return func(format string, v ...interface{}) {
		f := fmt.Sprintf("[%10s]: %s", context, format)
		log.Printf(f, v...)
	}
}
