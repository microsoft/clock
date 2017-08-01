package clock

import (
	"testing"
	"time"
)

const timeout = 10 * time.Millisecond

func assertGets(t *testing.T, c <-chan time.Time, fail string) {
	select {
	case <-c:
	case <-time.After(timeout):
		t.Error(fail)
	}
}

func assertDoesntGet(t *testing.T, c <-chan time.Time, fail string) {
	select {
	case now := <-c:
		t.Errorf("Failed at time %d: "+fail, now.UnixNano())
	case <-time.After(timeout):
	}
}

func assertBool(t *testing.T, expected, actual bool, fail string) {
	if expected != actual {
		t.Errorf("Expected %t to be %t:"+fail, actual, expected)
	}
}

func TestAfterGetsNegative(t *testing.T) {
	c := NewMockClock()
	assertGets(t, c.After(-time.Millisecond), "expected negative time to elapse immediately")
}

func TestAfterGetsExact(t *testing.T) {
	c := NewMockClock()
	ch := c.After(time.Millisecond * 2)
	assertDoesntGet(t, ch, "expected after not to get immediately")
	c.AddTime(time.Millisecond)
	assertDoesntGet(t, ch, "expected after not to get after 1 ms")
	c.AddTime(time.Millisecond)
	assertGets(t, ch, "expected after to get after 1 ms")
}

func TestAfterGetsOvershoot(t *testing.T) {
	c := NewMockClock()
	ch := c.After(time.Millisecond * 2)
	c.AddTime(time.Millisecond * 100)
	assertGets(t, ch, "expected after to get overshot")
}

func TestSleepWorks(t *testing.T) {
	// sleep just wraps .After internally, so just make sure it looks like it works
	c := NewMockClock()
	ch := make(chan time.Time)
	go func() {
		c.Sleep(10 * time.Millisecond)
		ch <- c.Now()
	}()

	assertDoesntGet(t, ch, "expected sleep to not get immediately")
	c.AddTime(time.Millisecond * 10)
	assertGets(t, ch, "expected sleep to return eventually")
}

func TestTickerWorks(t *testing.T) {
	c := NewMockClock()
	tk := c.NewTicker(5 * time.Millisecond)
	assertDoesntGet(t, tk.Chan(), "expected ticker to not get immediately")
	c.AddTime(4 * time.Millisecond)
	assertDoesntGet(t, tk.Chan(), "expected ticker to not get after 4ms")
	c.AddTime(5 * time.Millisecond)
	assertGets(t, tk.Chan(), "expected ticker to get after interval")
	c.AddTime(5 * time.Millisecond)
	assertGets(t, tk.Chan(), "expected ticker to keep getting")
	tk.Stop()
	c.AddTime(10 * time.Millisecond)
	assertDoesntGet(t, tk.Chan(), "expected ticker not to get after stopped")
}

func TestTickerDeadlock(t *testing.T) {
	c := NewMockClock()
	tk := c.NewTicker(5 * time.Millisecond)
	c.AddTime(6 * time.Millisecond)
	time.Sleep(1 * time.Millisecond)
	tk.Stop()
}

func TestTickerCatchesUp(t *testing.T) {
	c := NewMockClock()
	tk := c.NewTicker(5 * time.Millisecond)
	c.AddTime(20 * time.Millisecond)

	assertGets(t, tk.Chan(), "ticker sends catchup 1")
	assertGets(t, tk.Chan(), "ticker sends catchup 2")
	assertGets(t, tk.Chan(), "ticker sends catchup 3")
	assertGets(t, tk.Chan(), "ticker sends catchup 4")
	assertDoesntGet(t, tk.Chan(), "ticker catchup doesn't overshoot")
}

func TestTimerWorks(t *testing.T) {
	c := NewMockClock()
	tm := c.NewTimer(5 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer to not get immediately")
	c.AddTime(4 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer to not get after 4ms")
	c.AddTime(5 * time.Millisecond)
	assertGets(t, tm.Chan(), "expected timer to get after interval")
	c.AddTime(5 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer not to keep getting")
	assertBool(t, false, tm.Stop(), "expected reset to return false after timer finishes")
}

func TestTimerResets(t *testing.T) {
	c := NewMockClock()
	tm := c.NewTimer(5 * time.Millisecond)
	c.AddTime(4 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer to not get after 4ms")
	assertBool(t, true, tm.Reset(time.Millisecond*3), "expected reset to return true while timer running")
	c.AddTime(2 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer not to get after reset before interval is up")
	c.AddTime(time.Millisecond)
	assertGets(t, tm.Chan(), "expected timer to get after reset when interval is up")
	assertBool(t, false, tm.Stop(), "expected stop to return false after timer finishes")

	c.AddTime(100 * time.Millisecond)
	assertDoesntGet(t, tm.Chan(), "expected timer to not get after finished")

	assertBool(t, false, tm.Reset(time.Millisecond*3), "expected reset to return false after timer finishes")
	c.AddTime(3 * time.Millisecond)
	assertGets(t, tm.Chan(), "expected timer to get after reset when interval is up after restarted")
}
