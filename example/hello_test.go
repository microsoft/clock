package main

import (
	"testing"
	"time"

	"github.com/mixer/clock"
	"github.com/stretchr/testify/assert"
)

func TestDisplaysCorrectly(t *testing.T) {
	date, _ := time.Parse(time.UnixDate, "Sat Mar  7 11:12:39 PST 2015")
	c := clock.NewMockClock(date)
	d := displayer{c}

	assert.Equal(t, "11:12:39", d.formatted())
	c.AddTime(42 * time.Second)
	assert.Equal(t, "11:13:21", d.formatted())
}
