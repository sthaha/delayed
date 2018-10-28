// +build !debug

package delayed

func logger(context string) func(string, ...interface{}) {
	return func(_ string, _ ...interface{}) {}
}
