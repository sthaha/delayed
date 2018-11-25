package testutils

import (
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type Args map[string]interface{}

type Info struct {
	Called int
	Arg    Args
}

// Spy is used to inspect how many times hooks are called
type Spy struct {
	m      sync.RWMutex
	calls  int
	fnArgs Args
	Done   chan Info
}

func (s *Spy) Hook(v Args) func() {
	return func() { s.updateCalled(v) }
}

func (s *Spy) HookSync() (func(v Args) func(), chan Info) {
	ch := make(chan Info)

	hook := func(v Args) func() {
		return func() {
			s.updateCalled(v)
			select {
			case ch <- Info{s.calls, v}:
			default:
			}

		}
	}
	return hook, ch
}

func (s *Spy) updateCalled(v Args) {
	s.m.Lock()
	defer s.m.Unlock()
	s.calls++
	s.fnArgs = v
	Logger("spy")("called [%d] times with args: %s", s.calls, spew.Sdump(v))
}

func (s *Spy) Called() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.calls
}

func (s *Spy) Args() Args {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.fnArgs
}
