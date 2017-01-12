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
  // max size of any escape seq
	escBufferSize = 5

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
	// buffer to hold ESC characters
	escBuffer *bytes.Buffer
	// is input ready to be read
	ready bool
	// input in friendly form
	input string
	// if special keypress happened, store type here
	keypress int
	// marks for special char input
	mark bool
}

// InitStdin initializes our Stdin shim by hijacking input buffering
// and disabling echoing chars to the screen
func InitStdin() *Stdin {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	return &Stdin{
		escBuffer: bytes.NewBuffer(make([]byte, escBufferSize)),
		buffer:    bytes.NewBuffer(make([]byte, maxBufferSize)),
	}
}

// ReadNext readies the next line of input or keypress
func (s *Stdin) ReadNext() bool {
	s.keypress = None
	s.ready = false
	// read the byte off the wire
	b := make([]byte, 1)
	os.Stdin.Read(b)

	switch {
	case s.mark:
		s.escBuffer.Write(b)
		if isArrowUp(s.escBuffer.Bytes()) {
			s.keypress = ArrowUp
			s.mark = false
		} else if isArrowDown(s.escBuffer.Bytes()) {
			s.keypress = ArrowDown
			s.mark = false
		} else if s.escBuffer.Len() == escBufferSize {
			// didn't match on any special chars
			s.escBuffer.Reset()
			s.mark = false
		}
	case isEsc(b):
		// this could be the start of something special
		s.escBuffer.Write(b)
		s.mark = true
	case isEnter(b):
		s.flushBuffer()
		fmt.Printf(string(b))
	case isTab(b):
		s.keypress = Tab
	case s.buffer.Len() < maxBufferSize:
		s.buffer.Write(b)
		fmt.Print(string(b))
	default:
		// drop it like it's hot
	}

  // XXX: Always returns true, but maybe this is fine
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

func (s *Stdin) flushBuffer() {
	s.input = s.buffer.String()
	s.ready = true
	s.buffer.Reset()
}

func isEnter(b []byte) bool {
	return isEqual(b, []byte{10})
}

func isTab(b []byte) bool {
	return isEqual(b, []byte{9})
}

func isEsc(b []byte) bool {
	return isEqual(b, []byte{27})
}

func isArrowUp(b []byte) bool {
	return isEqual(b, []byte{27, 91, 65})
}

func isArrowDown(b []byte) bool {
	return isEqual(b, []byte{27, 91, 66})
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
