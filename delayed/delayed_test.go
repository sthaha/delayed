package delayed

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type args map[string]interface{}

type spy struct {
	calls  int
	fnArgs args

	t *testing.T
	m *sync.RWMutex
}

func newSpy(t *testing.T) *spy {
	return &spy{t: t, m: &sync.RWMutex{}}
}

func (s *spy) hook(v args) func() {
	return func() {
		s.m.Lock()
		defer s.m.Unlock()
		s.calls++
		s.fnArgs = v
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

func Test_simple(t *testing.T) {
	s := newSpy(t)

	df := &Fn{}
	df.Call(200*time.Millisecond, s.hook(args{"version": 1}))

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, 1, s.args()["version"])
}

func Test_delayed_orig(t *testing.T) {

	df := &Fn{}
	df.Call(3*time.Second, func() {
		log.Println("fn 1: done")
	})

	log.Println("waiting for 2 seconds")
	time.Sleep(2 * time.Second)

	log.Println("delayed after 5 seconds")

	df.Call(5*time.Second, func() {
		log.Println("fn 2: done")
	})

}

func Test_delayed_cancel(t *testing.T) {
	s := newSpy(t)
	df := &Fn{}
	df.Call(300*time.Millisecond, s.hook(args{"version": 1}))
	time.Sleep(50 * time.Millisecond)
	df.Call(300*time.Millisecond, s.hook(args{"version": 2}))
	time.Sleep(280 * time.Millisecond)
	df.Call(100*time.Millisecond, s.hook(args{"version": 3}))
	time.Sleep(120 * time.Millisecond)

	assert.Equal(t, 1, s.called, "must be called once")
	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, 1, s.called, "must be called once")
	assert.Equal(t, 3, s.args()["version"])
}
