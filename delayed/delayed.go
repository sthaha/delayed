package delayed

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Fn represents a fun to called after some duration has elapsed
type Fn struct {
	t *time.Timer
	m sync.Mutex
}

func logger(context string) func(string, ...interface{}) {
	return func(format string, v ...interface{}) {
		f := fmt.Sprintf("[%10s]: %s", context, format)
		log.Printf(f, v...)
	}
}

// Call waits for the duration to elapse and then calls f
// in its own goroutine. It returns a DelayedCall that can be used to
// over
func (df *Fn) Call(d time.Duration, fn func()) *Fn {
	df.m.Lock()
	defer df.m.Unlock()

	info := logger("delayed")
	if df.t != nil {
		info("stopping old")
		df.t.Stop()
	}

	info("Scheduled to run after %v", d)
	df.t = time.AfterFunc(d, fn)
	return df
}
