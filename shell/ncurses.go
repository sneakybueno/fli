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
	Delete
	Tab
)

// Stdin functions like normal stdin but allows us to catch keypress events
type Stdin struct {
	// buffer to hold input
	buffer *bytes.Buffer
	// buffer to hold ESC characters
	escBuffer *bytes.Buffer
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
		escBuffer: bytes.NewBuffer([]byte{}),
		buffer:    bytes.NewBuffer([]byte{}),
	}
}

// ReadNext appends to the input buffer and/or handles keypresses appropriately
func (s *Stdin) ReadNext() bool {
	// read the byte off the wire
	b := make([]byte, 1)
	if _, err := os.Stdin.Read(b); err != nil {
		// TODO: Handle this better
		return false
	}

	s.keypress = None
	switch {
	case s.mark:
		s.escBuffer.Write(b)
		switch {
		case isArrowUp(s.escBuffer.Bytes()):
			s.keypress = ArrowUp
			s.resetEscBuffer()
		case isArrowDown(s.escBuffer.Bytes()):
			s.keypress = ArrowDown
			s.resetEscBuffer()
		case s.escBuffer.Len() == escBufferSize:
			// didn't match on any special chars (prob not supported)
			s.resetEscBuffer()
		}
	case isEsc(b):
		// this could be the start of something special
		s.mark = true
	case isEnter(b):
		s.keypress = Enter
		fmt.Print(string(b))
	case isDel(b):
		s.keypress = Delete
		if s.buffer.Len() != 0 {
			s.buffer.Truncate(s.buffer.Len() - 1)
		}
	case isTab(b):
		s.keypress = Tab
	case s.buffer.Len() < maxBufferSize:
		s.buffer.Write(b)
		fmt.Print(string(b))
	default:
		// drop it like it's hot
	}

	return true
}

// Keypress returns type of keypress event
func (s *Stdin) KeyPress() int {
	return s.keypress
}

// Peek at the contents of the input buffer without destroying them
func (s *Stdin) Peek() string {
	return s.buffer.String()
}

// Flush the current input buffer
func (s *Stdin) Flush() string {
	return s.flushBuffer()
}

// Override replaces the input buffer contents with the string passed in
// Note: Not crazy about this (things should only flow up from this level, not down)
// but best way to handle shell history currently
func (s *Stdin) Override(input string) error {
	s.buffer.Reset()
	_, err := s.buffer.WriteString(input)
	return err
}

// Restore normal operations. Must call this before the program exits.
func (s *Stdin) Restore() {
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func (s *Stdin) flushBuffer() string {
	input := s.buffer.String()
	s.buffer.Reset()
	return input
}

func (s *Stdin) resetEscBuffer() {
	s.escBuffer.Reset()
	s.mark = false
}

// ASCII Matching
// ----------------------------------------------------------------------------

func isEnter(b []byte) bool {
	return isEqual(b, []byte{10})
}

func isDel(b []byte) bool {
	return isEqual(b, []byte{127})
}

func isTab(b []byte) bool {
	return isEqual(b, []byte{9})
}

func isEsc(b []byte) bool {
	return isEqual(b, []byte{27})
}

func isArrowUp(b []byte) bool {
	return isEqual(b, []byte{91, 65})
}

func isArrowDown(b []byte) bool {
	return isEqual(b, []byte{91, 66})
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
