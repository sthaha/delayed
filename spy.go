package delayed

import (
	"sync"

	"github.com/davecgh/go-spew/spew"
)

// spy is used to inspect how many times hooks are called
type args map[string]interface{}

type spy struct {
	m      sync.RWMutex
	calls  int
	fnArgs args
}

func (s *spy) hook(v args) func() {
	return func() {
		s.m.Lock()
		defer s.m.Unlock()
		s.calls++
		s.fnArgs = v
		logger("spy")("called [%d] times with args: %s", s.calls, spew.Sdump(v))
	}
}

func (s *spy) called() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.calls
}

func (s *spy) args() args {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.fnArgs
}
