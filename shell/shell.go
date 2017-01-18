package shell

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkg/term"
)

var (
	backspaceBytes = []byte("\b \b")
	newLineBytes   = []byte("\n")
)

// Shell struct for keeping track of shell things
type Shell struct {
	term *term.Term

	// buffer to hold input
	buffer *bytes.Buffer

	history *CmdHistory

	input string
	err   error
}

// Init creates a shell-like env
func Init() (*Shell, error) {
	t, err := term.Open("/dev/tty")
	if err != nil {
		return nil, err
	}

	return &Shell{
		term:    t,
		buffer:  bytes.NewBuffer([]byte{}),
		history: InitCmdHistory(50),
	}, nil
}

// Getters
// ----------------------------------------------------------------------------

// Input returns any available input
func (s *Shell) Input() string {
	return s.input
}

func (s *Shell) Error() error {
	return s.err
}

// Next returns true if the enter key has been pressed
func (s *Shell) Next() bool {
	for {
		c, err := s.getchar()
		if err != nil {
			s.err = err
			return false
		}

		switch {
		case isEnter(c):
			s.input = s.flushBuffer()
			if nonEmpty(s.input) {
				s.history.Add(s.input)
			}

			s.term.Write(newLineBytes)

			return true
		case isDelete(c):
			if s.buffer.Len() > 0 {
				s.buffer.Truncate(s.buffer.Len() - 1)
				s.term.Write(backspaceBytes)
			}
		case isArrowUp(c):
			// Need to stop when we reach position 0
			previousInput := s.history.Prev()
			s.overwriteBufferOnScreen(previousInput)
		case isArrowDown(c):
			// Need to stop when we reach position last position
			nextInput := s.history.Next()
			s.overwriteBufferOnScreen(nextInput)

		case isArrowLeft(c):
		case isArrowRight(c):
		case isCtrlC(c):
			fmt.Print("Closing app")
			s.Cleanup()
			os.Exit(0)
		default:
			s.buffer.Write(c)
			s.term.Write(c)
		}
	}
}

func (s *Shell) getchar() ([]byte, error) {
	s.term.SetRaw()

	// not sure if this should be hardcoded as 3 chars
	bytes := make([]byte, 3)
	numRead, err := s.term.Read(bytes)
	if err != nil {
		return nil, err
	}

	s.term.Restore()

	return bytes[0:numRead], nil
}

func (s *Shell) flushBuffer() string {
	input := s.buffer.String()
	s.buffer.Reset()

	return input
}

func (s *Shell) overwriteBufferOnScreen(buffer string) {
	bufferBytes := []byte(buffer)

	// delete everything in current buffer
	length := s.buffer.Len()
	for i := 0; i < length; i++ {
		s.term.Write(backspaceBytes)
	}

	s.flushBuffer()

	s.buffer.Write(bufferBytes)
	s.term.Write(bufferBytes)
}

// Cleanup does any work needed to cleanly close the shell
func (s *Shell) Cleanup() {
	s.term.Restore()
	s.term.Close()
}
