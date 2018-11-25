package delayed_test

import (
	"sync"
	"time"

	"github.com/sthaha/delayed"
	"github.com/sthaha/delayed/testutils"
)

type idlerstat int

const (
	running idlerstat = 0
	idled   idlerstat = 1
)

var stat = running

func status() idlerstat {
	return stat
}

var unidleRequested = false

var fn delayed.Fn

func unidle() {
	info := testutils.Logger("unidle")

	info("-------------------------------")
	if unidleRequested {
		info("Already scheduled to unidle, ignoring this call")
		return
	}

	unidleRequested = true

	info("unidle: going to delayed unidle")

	fn.Reset(3*time.Second, func() {
		info("unidle: unidling ...")
		if status() == running {
			info("unidle: No need to idle since it is already unidled")
			return
		}
		stat = running
		unidleRequested = false
	})
}

func idle() {
	info := testutils.Logger("idle")

	info("-------------------------------")
	if unidleRequested {
		info("Ignoring Idle request since unidle has been requested earlier")
		return
	}

	info("going to delayed idle")

	fn.Reset(3*time.Second, func() {

		if status() == idled {
			info("No need to idle since it is already idled")
			return
		}
		info("idling ... - START")
		time.Sleep(500 * time.Millisecond)
		info("idling ... - DONE")
		stat = idled
	})
}

func simulateIdler(wg *sync.WaitGroup) {
	defer wg.Done()

	idle()
	time.Sleep(200 * time.Millisecond)
	idle()
	time.Sleep(200 * time.Millisecond)
	unidle()
	time.Sleep(200 * time.Millisecond)
	idle()
	time.Sleep(200 * time.Millisecond)
	unidle()
	time.Sleep(5 * time.Second)
	idle()
}

func Example_idler() {
	info := testutils.Logger("main")
	wg := &sync.WaitGroup{}

	wg.Add(1)
	info("starting test run")
	simulateIdler(wg)
	wg.Wait()
	info("done")
	// HACK: output below is used as hack to run the test
	// run this using go test -race -run Example

	// Output:
	//
}
