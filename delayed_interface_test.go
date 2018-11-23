package delayed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type caller interface {
	Call() error
}

type canceller interface {
	Cancel() bool
}

type resetter interface {
	Reset(time.Duration, func()) error
	ResetDelay(time.Duration) error
	ResetFunc(func()) error
}

type delayedFn interface {
	caller
	canceller
	resetter
}

func TestInterface(t *testing.T) {
	var c caller = &Fn{}
	assert.Error(t, c.Call())
}

func TestInterface_caller(t *testing.T) {
	var c caller = &Fn{}
	assert.Error(t, c.Call())
}

func TestInterface_canceller(t *testing.T) {
	var c canceller = &Fn{}
	assert.False(t, c.Cancel())
}

func TestInterface_resetter(t *testing.T) {
	var r resetter = &Fn{}
	assert.Error(t, r.Reset(-1, func() {}))
	assert.Error(t, r.ResetDelay(-1))
	var f func()
	assert.Error(t, r.ResetFunc(f))
}

func TestInterface_works(t *testing.T) {
	var c delayedFn = &Fn{}
	assert.Error(t, c.Call())

	s := &spy{}

	err := c.Reset(20*time.Millisecond, s.hook(args{"version": 1}))
	assert.NoError(t, err, "must be created")

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, 1, s.args()["version"])
	assert.False(t, c.Cancel(), "cancelling already executed must return false")
}
