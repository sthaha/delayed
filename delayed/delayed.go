package delayed

import (
	"sync"
	"time"
)

// Fn represents a fun to called after some duration has elapsed
type Fn struct {
	t *time.Timer
	m sync.Mutex
}

var debug = logger("delayed")

// Call waits for the duration to elapse and then calls f
// in its own goroutine. It returns a DelayedCall that can be used to
// over
func (f *Fn) Call(d time.Duration, fn func()) *Fn {

	f.Cancel()

	debug("Scheduled to run after %v", d)

	f.m.Lock()
	defer f.m.Unlock()
	f.t = time.AfterFunc(d, fn)
	return f
}

// Cancel cancels the function that was scheduled by Call
func (f *Fn) Cancel() bool {
	f.m.Lock()
	defer f.m.Unlock()

	if f.t == nil {
		return false
	}

	debug("cancelling delayed call")
	f.t.Stop()
	return true
}
