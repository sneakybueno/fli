package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const (
	// max size of stdin buffer
	maxBufferSize = 1024
	specialSeqSize = 3

	// "enums" for keypresses
	None = iota
	Enter
	ArrowUp
	ArrowDown
	Tab
)

// Stdin functions like normal stdin but allows us to catch keypress events
type Stdin struct {
	// buffer to hold input
	buffer *bytes.Buffer
	// number of bytes currently written to buffer
	written int
	// is input ready to be read
	ready	 bool
	// input in friendly form
	input  string
	// if special keypress happened, store type here
	keypress    int
}

// InitStdin initializes our Stdin shim by hijacking input buffering
// and disabling echoing chars to the screen
func InitStdin() *Stdin {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	return &Stdin{
		buffer: bytes.NewBuffer(make([]byte, maxBufferSize)),
	}
}

// ReadNext readies the next line of input or keypress
// XXX: Always returns true, but maybe this is fine
// TODO: Refactor this
func (s *Stdin) ReadNext() bool {
	s.keypress = None
	s.ready = false
	// read the byte off the wire
	b := make([]byte, 1)
	os.Stdin.Read(b)

	if isMulti(b) {
		// this could be the start of something special
		bspecial := []byte{}
		bspecial = append(bspecial, b...)
		// read twice more
		for i := 0; i < specialSeqSize-1; i++ {
			os.Stdin.Read(b)
			bspecial = append(bspecial, b...)
		}
		if isArrowUp(bspecial) {
			s.keypress = ArrowUp
		} else if isArrowDown(bspecial) {
			s.keypress = ArrowDown
		}
	} else if isEnter(b) {
		s.input = s.buffer.String()
		s.ready = true
		s.written = 0
		s.buffer.Reset()
		fmt.Printf(string(b))
	} else if isTab(b) {
		s.keypress = Tab
	} else if s.written < maxBufferSize {
		s.buffer.Write(b)
		s.written++
		fmt.Print(string(b))
	} else {
		// drop it like it's hot
	}

	return true
}

// Keypress returns type of keypress event
func (s *Stdin) KeyPress() int {
	return s.keypress
}

// Text returns input string if it is ready (ie Enter has been pressed)
func (s *Stdin) Text() (string, bool) {
	isReady := s.ready
	s.ready = false
	return s.input, isReady
}

// Restore normal operations. Must call this before the program exits.
func (s *Stdin) Cleanup() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func isEnter(b []byte) bool {
	return isEqual(b, []byte{10})
}

func isTab(b []byte) bool {
	return isEqual(b, []byte{9})
}

// XXX: need to do more research but figure out if byte 27 is special
func isMulti(b []byte) bool {
	return isEqual(b, []byte{27})
}

func isArrowUp(b []byte) bool {
	return isEqual(b, []byte{27,91,65})
}

func isArrowDown(b []byte) bool {
	return isEqual(b, []byte{27,91,66})
}

func isEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range b {
		if a[i] != c {
			return false
		}
	}
	return true
}
