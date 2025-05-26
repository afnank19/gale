package main

import (
	"fmt"
	"time"
)

type State struct {
	workers     int
	concurrency int
	duration    time.Duration
}

func main() {
	fmt.Println("GALE! Hammer your server")

	// ParseArgs()
	timer1 := time.NewTimer(2 * time.Second)

	// The `<-timer1.C` blocks on the timer's channel `C`
	// until it sends a value indicating that the timer
	// fired.
	fmt.Println("Timer 1 fired", <-timer1.C)
}
