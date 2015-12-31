package clock

import (
	"time"
)

// DefaultClock is an implementation of the Clock interface that uses standard
// time methods.
type DefaultClock struct{}

func (dc DefaultClock) Now() time.Time                         { return time.Now() }
func (dc DefaultClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (dc DefaultClock) Sleep(d time.Duration)                  { time.Sleep(d) }
func (dc DefaultClock) Tick(d time.Duration) <-chan time.Time  { return time.Tick(d) }
func (dc DefaultClock) AfterFunc(d time.Duration, f func()) Timer {
	return &defaultTimer{*time.AfterFunc(d, f)}
}
func (dc DefaultClock) NewTimer(d time.Duration) Timer {
	return &defaultTimer{*time.NewTimer(d)}
}
func (dc DefaultClock) NewTicker(d time.Duration) Ticker {
	return &defaultTicker{*time.NewTicker(d)}
}

type defaultTimer struct{ time.Timer }

var _ Timer = new(defaultTimer)

func (d *defaultTimer) Chan() <-chan time.Time {
	return d.C
}

type defaultTicker struct{ time.Ticker }

var _ Ticker = new(defaultTicker)

func (d *defaultTicker) Chan() <-chan time.Time {
	return d.C
}

// Default clock that uses time.Now as its time source. This is what you should
// normally use in your code.
var C = DefaultClock{}
