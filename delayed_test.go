package delayed

import (
	"testing"
	"time"

	"github.com/sthaha/delayed/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewFn(t *testing.T) {
	f := NewFn(5*time.Second, func() {})
	assert.True(t, f.valid())
	assert.False(t, f.Cancel())
}

func TestNewFn_invalid_duration(t *testing.T) {
	f := NewFn(-2*time.Second, func() {})
	assert.False(t, f.valid())
	assert.Error(t, f.Call(), "call to invalid Fn will return error")
}

func TestNewFn_invalid_fn(t *testing.T) {
	f := &Fn{d: 10 * time.Millisecond}
	assert.False(t, f.valid())
	assert.Error(t, f.Call(), "call to invalid Fn will return error")
}

func Test_Call(t *testing.T) {
	s := &testutils.Spy{}

	fn, err := Call(20*time.Millisecond, s.Hook(testutils.Args{"version": 1}))
	assert.NoError(t, err, "must be created")

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")
	assert.Equal(t, 1, s.Args()["version"])
	assert.False(t, fn.Cancel(), "cancelling already executed must return false")
}

func Test_Call_error(t *testing.T) {
	s := &testutils.Spy{}

	validFn := func() {}
	_, err := Call(-1*time.Millisecond, validFn)
	assert.Error(t, err, "invoking Call with invalid args must return error")

	var invalidFn func()

	_, err = Call(10*time.Millisecond, invalidFn)
	assert.Error(t, err, "invoking Call with invalid args must return error")

	_, err = Call(10*time.Millisecond, s.Hook(testutils.Args{"v": 1}))
	assert.NoError(t, err)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")

	time.Sleep(30 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")
	assert.Equal(t, 1, s.Args()["v"])
}

func TestEmpty_call(t *testing.T) {
	fn := &Fn{}
	assert.Error(t, fn.Call(), "invalid fn did not return error")

	// Resetting Delay will still
	assert.Error(t, fn.ResetDelay(0), "invalid duration must return error")
	assert.Error(t, fn.ResetDelay(2*time.Millisecond), "invalid duration must return error")
}

func TestFn_call(t *testing.T) {
	fn := &Fn{}

	s := &testutils.Spy{}
	fn.Reset(20*time.Millisecond, s.Hook(testutils.Args{"version": 1}))

	time.Sleep(25 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")

	time.Sleep(25 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")
	assert.Equal(t, 1, s.Args()["version"])
}

func Test_simple_cancel(t *testing.T) {

	fn := &Fn{}
	assert.False(t, fn.Cancel(), "cancel on unscheduled must return false")

	s := &testutils.Spy{}
	fn.Reset(200*time.Millisecond, s.Hook(testutils.Args{"version": 1}))
	time.Sleep(100 * time.Millisecond)

	assert.True(t, fn.Cancel(), "cancel on scheduled must return true")
	assert.Equal(t, 0, s.Called(), "must be called once")

	// ensure it is not called
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 0, s.Called(), "must be called once")
}

func Test_multi(t *testing.T) {
	s := &testutils.Spy{}

	fn := &Fn{}

	const n = 5
	for i := 1; i <= n; i++ {

		go func(i int) {
			testutils.Logger("test")("calling with %v", i)
			delay := time.Duration(20*i) * time.Millisecond

			// access to call must not cause any race
			fn.Reset(delay, s.Hook(testutils.Args{"version": i}))
		}(i)
		time.Sleep(15 * time.Millisecond)
	}

	// should execute in 75 (15 * 5) millisecond
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, s.Called(), "must be called once")
	assert.Equal(t, n, s.Args()["version"])

	// ensure it is not called again after some time
	time.Sleep(250 * time.Millisecond)
	assert.Equal(t, 1, s.Called(), "must be called once")
}

func Test_cancel_calls(t *testing.T) {
	s := &testutils.Spy{}
	fn := &Fn{}

	fn.Reset(300*time.Millisecond, s.Hook(testutils.Args{"version": 1}))
	time.Sleep(50 * time.Millisecond)
	fn.Reset(300*time.Millisecond, s.Hook(testutils.Args{"version": 2}))
	time.Sleep(280 * time.Millisecond)
	fn.Reset(100*time.Millisecond, s.Hook(testutils.Args{"version": 3}))
	time.Sleep(120 * time.Millisecond)

	assert.Equal(t, 1, s.Called(), "must be called once")
	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, 1, s.Called(), "must be called once")
	assert.Equal(t, 3, s.Args()["version"])
}
