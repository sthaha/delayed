package delayed_test

import (
	"testing"
	"time"

	"github.com/sthaha/delayed"
	"github.com/sthaha/delayed/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFakeFn_is_callable(t *testing.T) {
	spy := &testutils.Spy{}

	hook, wait := spy.HookSync()
	defer close(wait)

	// test that FakeFn satisfies Callable
	var fn delayed.Callable = delayed.NewFakeFn(
		2000*time.Millisecond,
		hook(testutils.Args{"v": 1}),
	)

	m, ok := fn.(*delayed.FakeFn)

	m.Call()
	m.Clock.Advance(2000 * time.Millisecond)
	c := <-wait

	assert.True(t, ok, "failed type conversion to FakeFn")
	assert.Equal(t, 1, c.Called)
	assert.Equal(t, 1, c.Arg["v"])
}

func TestFake_cancel(t *testing.T) {
	spy := &testutils.Spy{}
	hook, wait := spy.HookSync()
	defer close(wait)

	var fn = delayed.FakeCall(
		2000*time.Millisecond,
		hook(testutils.Args{"v": 1}),
	)

	fn.Clock.Advance(1000 * time.Millisecond)
	cancelled := fn.Cancel()

	assert.True(t, cancelled, "must not have run")
	assert.Equal(t, 0, spy.Called())
	assert.Nil(t, spy.Args()["v"])
}

func TestFake_reset(t *testing.T) {

	spy := &testutils.Spy{}
	hook, wait := spy.HookSync()
	defer close(wait)

	fn := delayed.FakeCall(
		2000*time.Millisecond,
		hook(testutils.Args{"v": 1}),
	)

	fn.Clock.Advance(1000 * time.Millisecond)

	// reset
	fn.Reset(
		1000*time.Millisecond,
		hook(testutils.Args{"v": 2}),
	)

	fn.Clock.Advance(1000 * time.Millisecond)
	c := <-wait

	cancelled := fn.Cancel()
	assert.False(t, cancelled, "must not have run")
	assert.Equal(t, 1, c.Called)
	assert.Equal(t, 2, c.Arg["v"])
}

func TestFake_resetFn(t *testing.T) {

	spy := &testutils.Spy{}
	hook, wait := spy.HookSync()
	defer close(wait)

	fn := delayed.FakeCall(
		2000*time.Millisecond,
		hook(testutils.Args{"v": 1}),
	)

	fn.Clock.Advance(1000 * time.Millisecond)

	// reset
	fn.ResetFunc(hook(testutils.Args{"v": 2}))

	fn.Clock.Advance(2000 * time.Millisecond)
	c := <-wait

	cancelled := fn.Cancel()
	assert.False(t, cancelled, "must not have run")
	assert.Equal(t, 1, c.Called)
	assert.Equal(t, 2, c.Arg["v"])
}

func TestFake_resetDelay(t *testing.T) {

	spy := &testutils.Spy{}
	hook, wait := spy.HookSync()
	defer close(wait)

	fn := delayed.FakeCall(
		8000*time.Millisecond,
		hook(testutils.Args{"v": 1}),
	)

	fn.Clock.Advance(1000 * time.Millisecond)

	// reset
	fn.ResetDelay(100 * time.Millisecond)

	fn.Clock.Advance(2000 * time.Millisecond)
	c := <-wait

	cancelled := fn.Cancel()
	assert.False(t, cancelled, "must not have run")
	assert.Equal(t, 1, c.Called)
	assert.Equal(t, 1, c.Arg["v"])
}
