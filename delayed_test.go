package delayed

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Call(t *testing.T) {
	s := newSpy(t)

	fn := Call(200*time.Millisecond, s.hook(args{"version": 1}))

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, 1, s.args()["version"])
	assert.False(t, fn.Cancel())
}

func Test_fn_call(t *testing.T) {
	s := newSpy(t)

	fn := &Fn{}
	fn.Call(200*time.Millisecond, s.hook(args{"version": 1}))

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")

	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, 1, s.args()["version"])
}

func Test_simple_cancel(t *testing.T) {
	s := newSpy(t)

	fn := &Fn{}

	assert.False(t, fn.Cancel(), "cancel on unscheduled must return false")

	fn.Call(200*time.Millisecond, s.hook(args{"version": 1}))
	time.Sleep(100 * time.Millisecond)

	assert.True(t, fn.Cancel(), "cancel on scheduled must return true")
	assert.Equal(t, 0, s.called(), "must be called once")

	// ensure it is not called
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 0, s.called(), "must be called once")
}

func Test_multi(t *testing.T) {
	s := newSpy(t)

	fn := &Fn{}

	const n = 5
	for i := 1; i <= n; i++ {

		go func(i int) {
			logger("test")("calling with %v", i)
			delay := time.Duration(200*i) * time.Millisecond

			// access to call must not cause any race
			fn.Call(delay, s.hook(args{"version": i}))
		}(i)
		time.Sleep(150 * time.Millisecond)
	}

	// should execute in 750 (150 * 5) millisecond
	time.Sleep(1000 * time.Millisecond)

	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, n, s.args()["version"])

	// ensure it is not called again after some time
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.called(), "must be called once")
}

func Test_cancel_calls(t *testing.T) {
	s := newSpy(t)
	fn := &Fn{}

	fn.Call(300*time.Millisecond, s.hook(args{"version": 1}))
	time.Sleep(50 * time.Millisecond)
	fn.Call(300*time.Millisecond, s.hook(args{"version": 2}))
	time.Sleep(280 * time.Millisecond)
	fn.Call(100*time.Millisecond, s.hook(args{"version": 3}))
	time.Sleep(120 * time.Millisecond)

	assert.Equal(t, 1, s.called(), "must be called once")
	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, 1, s.called(), "must be called once")
	assert.Equal(t, 3, s.args()["version"])
}
