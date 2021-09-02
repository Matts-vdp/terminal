package progress

import (
	"fmt"

	"github.com/Matts-vdp/terminal/ter"
)

type progress struct {
	terminal *ter.TerminalOut
	line     int
	total    int
	current  int
	name     string
}

// Creates a new progress counter on line "line" of the given terminal.
// You can update the count by sending a new count on the returned chan.
// The initial string is returned
func InitProgresBar(terminal *ter.TerminalOut, line, total int, name string) chan int {
	ch := make(chan int)
	p := &progress{terminal: terminal, line: line, total: total, name: name}
	go func() {
		for new := range ch {
			p.update(new)
		}
	}()
	return ch
}

func (p *progress) string() string {
	return fmt.Sprintf("%d/%d %s", p.current, p.total, p.name)
}
func (p *progress) update(new int) {
	p.current = new
	p.terminal.UpdateLine(p.line, p.string())
}
