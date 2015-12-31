package clock

import (
	"fmt"
	"sync"
	"time"
)

// MockClock provides a Clock whose time only changes or advances when
// manually specified to do so. This is useful for unit tests.
type MockClock struct {
	// Mutex is used to guard the current time, and the cond is used to
	// broadcast changes to things waiting for the time to change.
	cond *sync.Cond
	now  time.Time
}

var _ Clock = new(MockClock)

func (m *MockClock) Now() time.Time {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	return m.now
}

func (m *MockClock) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	target := m.Now().Add(d)

	go func() {
		for {
			m.cond.L.Lock()
			if target.After(m.now) {
				m.cond.Wait()
			} else {
				m.cond.L.Unlock()
				ch <- m.now
				return
			}

			m.cond.L.Unlock()
		}
	}()

	return ch
}

func (m *MockClock) Sleep(d time.Duration) {
	<-m.After(d)
}

func (m *MockClock) Tick(d time.Duration) <-chan time.Time {
	return m.NewTicker(d).Chan()
}

func (m *MockClock) AfterFunc(d time.Duration, f func()) Timer {
	t := m.NewTimer(d)
	go func() {
		<-t.Chan()
		go f()
	}()

	return t
}

func (m *MockClock) NewTimer(d time.Duration) Timer {
	t := NewMockTimer(m)
	t.Reset(d)
	return t
}

func (m *MockClock) NewTicker(d time.Duration) Ticker {
	return NewMockTicker(m, d)
}

// Sets the mock clock's time to the given absolute time.
func (m *MockClock) SetTime(t time.Time) {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	assertFuture(m.now, t)
	m.now = t
	m.cond.Broadcast()
}

// Adds the given time duration to the clock.
func (m *MockClock) AddTime(d time.Duration) {
	m.cond.L.Lock()
	defer m.cond.L.Unlock()

	assertFuture(m.now, m.now.Add(d))
	m.now = m.now.Add(d)
	m.cond.Broadcast()
}

func assertFuture(a, b time.Time) {
	na, nb := a.UnixNano(), b.UnixNano()
	if na > nb {
		panic(fmt.Sprintf("Tried to tick backwards from %d to %d, but cannot travel into the past!", na, nb))
	}
}

// Creates a new mock clock, with its current time set to the provided
// optional start time.
func NewMockClock(start ...time.Time) *MockClock {
	m := &MockClock{cond: sync.NewCond(new(sync.Mutex))}

	if len(start) > 1 {
		panic(fmt.Sprintf("Expected one argument to clock.NewMock, got %d", len(start)))
	} else if len(start) == 1 {
		m.SetTime(start[0])
	} else {
		m.SetTime(time.Now())
	}

	return m
}
