package delayed

import (
	"fmt"
	"sync"
	"time"

	"github.com/sthaha/delayed/testutils"
)

// Callable is the interface a delayed function call must satisfy
type Callable interface {
	Call() error
	Cancel() bool
	Reset(time.Duration, func()) error
	ResetDelay(time.Duration) error
	ResetFunc(func()) error
}

// Fn represents a fun to called after some duration has elapsed
type Fn struct {
	m  sync.Mutex
	d  time.Duration
	fn func()

	t *time.Timer
}

var debug = testutils.Logger("delayed")

// NewFn returns an instance of delayed Fn
func NewFn(d time.Duration, fn func()) *Fn {
	return &Fn{d: d, fn: fn}
}

// Call executes a fn after the duration and returns a handle to the
// function so that it can be cancelled or overridden
func Call(d time.Duration, fn func()) (*Fn, error) {
	f := &Fn{d: d, fn: fn}
	if err := f.Call(); err != nil {
		return nil, err
	}

	return f, nil
}

// ResetFunc resets func to be invoked
func (f *Fn) ResetFunc(fn func()) error {
	f.m.Lock()
	defer f.m.Unlock()

	f.cancel()
	f.fn = fn
	return f.call()
}

// ResetDelay resets the delay and starts again
func (f *Fn) ResetDelay(d time.Duration) error {
	f.m.Lock()
	defer f.m.Unlock()

	f.cancel()

	debug("scheduled to run after %v", d)
	f.d = d
	return f.call()
}

// Reset resets both duration and the fn to call.
func (f *Fn) Reset(d time.Duration, fn func()) error {
	f.m.Lock()
	defer f.m.Unlock()

	f.cancel()
	debug("scheduled to run after %v", f.d)
	f.d = d
	f.fn = fn
	return f.call()
}

// Call waits for the duration to elapse and then calls fn in its own goroutine.
func (f *Fn) Call() error {
	f.m.Lock()
	defer f.m.Unlock()
	return f.call()
}

// Cancel cancels the function that was scheduled by Call
func (f *Fn) Cancel() bool {
	f.m.Lock()
	defer f.m.Unlock()
	return f.cancel()
}

func (f *Fn) call() error {

	if !f.valid() {
		return fmt.Errorf("invalid delayed function")
	}

	f.cancel()

	debug("Scheduled to run after %v", f.d)
	f.t = time.AfterFunc(f.d, f.fn)
	return nil
}

func (f *Fn) cancel() bool {
	if f.t == nil {
		return false
	}

	debug("cancelling delayed call")
	return f.t.Stop()

}

func (f *Fn) valid() bool {
	return f.fn != nil && f.d >= 0
}
