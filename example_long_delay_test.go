package delayed_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sthaha/delayed"
	"github.com/sthaha/delayed/testutils"
)

func Example_long_running_fn() {
	//
	debug := testutils.Logger("main")

	var count uint64

	fn, _ := delayed.Call(100*time.Millisecond, func() {
		debug := testutils.Logger("first")

		debug("going to execute")
		atomic.AddUint64(&count, 1)
		// simulate some activity while a Reset could be called from
		// else where in the code
		time.Sleep(500 * time.Millisecond)

		debug("execute done")
	})

	// give enough time for the call to run
	time.Sleep(110 * time.Millisecond)
	debug("going to schedule another")

	fn.Reset(100*time.Millisecond, func() {
		debug := testutils.Logger("second")

		debug("going to execute")
		atomic.AddUint64(&count, 1)

		time.Sleep(200 * time.Millisecond)
		debug("execute done")
	})

	time.Sleep(110 * time.Millisecond)

	called := atomic.LoadUint64(&count)
	fmt.Printf("not mutex protected thus count is %d\n", called)

	// wait for both to finish
	time.Sleep(800 * time.Millisecond)

	called = atomic.LoadUint64(&count)
	fmt.Printf("expect calls to be %d\n", called)

	// Output:
	// not mutex protected thus count is 2
	// expect calls to be 2
}

func Example_long_running_fn_fixed() {
	debug := testutils.Logger("main")

	var count uint64

	var m sync.Mutex

	fn, _ := delayed.Call(100*time.Millisecond, func() {
		m.Lock()
		defer m.Unlock()

		debug := testutils.Logger("first")

		debug("going to execute")
		atomic.AddUint64(&count, 1)
		// simulate some activity while a Reset could be called from
		// else where in the code
		time.Sleep(500 * time.Millisecond)

		debug("execute done")
	})

	// give enough time for the call to run
	time.Sleep(110 * time.Millisecond)
	debug("going to schedule another")

	fn.Reset(100*time.Millisecond, func() {
		m.Lock()
		defer m.Unlock()
		debug := testutils.Logger("second")

		debug("going to execute")
		atomic.AddUint64(&count, 1)

		time.Sleep(200 * time.Millisecond)
		debug("execute done")
	})

	time.Sleep(110 * time.Millisecond)

	called := atomic.LoadUint64(&count)
	fmt.Printf("mutex protected thus count is %d\n", called)

	// wait for both to finish
	time.Sleep(800 * time.Millisecond)

	called = atomic.LoadUint64(&count)
	fmt.Printf("expect calls to be %d\n", called)
	time.Sleep(800 * time.Millisecond)

	// Output:
	// mutex protected thus count is 1
	// expect calls to be 2
}
