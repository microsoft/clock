package main

import (
	"fmt"

	"github.com/WatchBeam/clock"
)

func main() {
	fmt.Printf("the time is %s", displayer{clock.C}.formatted())
}

type displayer struct {
	c clock.Clock
}

func (d displayer) formatted() string {
	now := d.c.Now()
	return fmt.Sprintf("%d:%d:%d", now.Hour(), now.Minute(), now.Second())
}
