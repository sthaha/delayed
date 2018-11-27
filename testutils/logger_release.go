// +build !debug

package testutils

func Logger(context string) func(string, ...interface{}) {
	return func(_ string, _ ...interface{}) {}
}
