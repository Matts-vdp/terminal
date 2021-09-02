package ter

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	cursorUp   = "\033[A"
	cursorDown = "\033[B"
)

type lineUpdate struct {
	pos int
	str string
}

type TerminalOut struct {
	maxLength  int
	lines      []string
	output     io.Writer
	autoUpdate bool
	updateChan chan lineUpdate
	doneChan   chan bool
	first      bool
}

// Sets the command prompt to virtual mode.
// This way escape characters can be used to control the cursor position
func setupVirtualTerminal() {
	const Vir = 0x4
	var mode uint32
	h := os.Stdout.Fd()
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procGetConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode)))
	procSetConsoleMode.Call(h, uintptr(mode|Vir))
}

// Create a terminalOut object that can be used to display and change lines
func InitTerminal(o io.Writer, lines []string, autoUpdate bool) *TerminalOut {
	setupVirtualTerminal()
	t := &TerminalOut{0, lines, o, autoUpdate, make(chan lineUpdate), make(chan bool), true}
	t.updateMax()
	go t.startLineUpdater()
	return t
}

func (t *TerminalOut) startLineUpdater() {
	for upd := range t.updateChan {
		t.lines[upd.pos] = upd.str
		t.updateMax()
		if t.autoUpdate {
			t.display()
		}
	}
	t.doneChan <- true
	close(t.doneChan)
}

// Calculates how big the largest passed string is
func (t *TerminalOut) updateMax() {
	max := t.maxLength
	for _, l := range t.lines {
		if size := len(l); size > max {
			max = size
		}
	}
	t.maxLength = max
}

func (t *TerminalOut) writeLines() {
	for _, l := range t.lines {
		f := "%-" + fmt.Sprint(t.maxLength) + "s\n"
		fmt.Fprintf(t.output, f, l)
	}
}

func (t *TerminalOut) toTop() {
	for i := 0; i < len(t.lines); i++ {
		fmt.Fprint(t.output, cursorUp)
	}
}

// Updates the displayed lines.
// Use this only for the first time or when autoupdate is false
func (t *TerminalOut) ManualRefresh() {
	if t.autoUpdate && !t.first {
		return
	}
	t.display()
}

func (t *TerminalOut) display() {
	if !t.first {
		t.toTop()
	} else {
		t.first = false
	}
	t.writeLines()
}

// Update the text on line "position" (starts at 0).
// Can used concurrently without problems
func (t *TerminalOut) UpdateLine(position int, newStr string) {
	t.updateChan <- lineUpdate{position, newStr}
}

// Stops the auto update and makes sure that all updates are displayed
func (t *TerminalOut) Close() {
	close(t.updateChan)
	<-t.doneChan
}
