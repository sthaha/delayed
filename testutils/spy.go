package testutils

import (
	"sync"

	"github.com/davecgh/go-spew/spew"
)

type Args map[string]interface{}

// Spy is used to inspect how many times hooks are called
type Spy struct {
	m      sync.RWMutex
	calls  int
	fnArgs Args
}

func (s *Spy) Hook(v Args) func() {
	return func() {
		s.m.Lock()
		defer s.m.Unlock()
		s.calls++
		s.fnArgs = v
		Logger("spy")("called [%d] times with args: %s", s.calls, spew.Sdump(v))
	}
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
