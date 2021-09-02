# terminal
A go library to display multiple lines in the terminal and update them.

For now it only works for windows

## Example
This example shows how to use the terminal and progress modules to display 3 progress counters.
```Go
package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/Matts-vdp/terminal/progress"
	"github.com/Matts-vdp/terminal/ter"
)

func count(c chan int, done chan bool) {
	for i := 0; i < 101; i++ {
		c <- i
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
	done <- true
}

func main() {
	l := make([]string, 3)
	t := ter.InitTerminal(os.Stdout, l, true)
	defer t.Close()
	t.ManualRefresh()
	done := make(chan bool)
	p1 := progress.InitProgresBar(t, 0, 100, "Progress 1")
	go count(p1, done)
	p2 := progress.InitProgresBar(t, 1, 100, "Progress 2")
	go count(p2, done)
	p3 := progress.InitProgresBar(t, 2, 100, "Progress 3")
	go count(p3, done)
	for i := 0; i < 3; i++ {
		<-done
	}
}
```
