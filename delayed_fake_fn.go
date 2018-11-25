package delayed

import (
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/sthaha/delayed/testutils"
)

// FakeFn satisfies a Callable where the clock can be controlled/mocked
type FakeFn struct {
	Clock clockwork.FakeClock

	d  time.Duration
	fn func()

	cancel chan bool
}

func NewFakeFn(d time.Duration, fn func()) *FakeFn {
	return &FakeFn{
		Clock:  clockwork.NewFakeClockAt(time.Now()),
		d:      d,
		fn:     fn,
		cancel: make(chan bool),
	}
}

func FakeCall(d time.Duration, fn func()) *FakeFn {
	f := NewFakeFn(d, fn)
	f.Call()
	return f
}

func (m *FakeFn) Call() error {
	log := testutils.Logger("call")

	go func() {
		log("waiting ....... %v to run func\n", m.d)
		select {
		case <-m.Clock.After(m.d):
			log("running .......")
			m.fn()
			log("running ....... done")
		case <-m.cancel:
			log("got cancel .......")
			return
		}
		log("waiting ....... done")
	}()
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (m *FakeFn) Cancel() bool {
	select {
	case m.cancel <- true:
		return true
	default:
		return false
	}
}

func (m *FakeFn) Reset(d time.Duration, fn func()) error {
	m.Cancel()
	m.d = d
	m.fn = fn
	m.Call()
	return nil
}

func (m *FakeFn) ResetDelay(d time.Duration) error {
	return m.Reset(d, m.fn)
}

func (m *FakeFn) ResetFunc(fn func()) error {
	return m.Reset(m.d, fn)
}
