package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"

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

	commands Commands
	history  *CmdHistory

	prompt string
	input  string
	err    error
}

// Init creates a shell-like env
func Init(prompt string) (*Shell, error) {
	t, err := term.Open("/dev/tty")
	if err != nil {
		return nil, err
	}

	s := &Shell{
		term:    t,
		buffer:  bytes.NewBuffer([]byte{}),
		history: InitCmdHistory(50),
	}

	s.prompt = prompt
	fmt.Print(s.prompt)

	s.AddCommand("exit", exitHandler)

	return s, nil
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

func (s *Shell) Prompt() string {
	return s.prompt
}

func (s *Shell) SetPrompt(prompt string) {
	s.prompt = prompt
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
		case isTab(c):
			// searchTerm := s.getLastWordFromBuffer()

			// fakeDatasource := []string{"hello", "users"}
			// tabCompletion, err := FindNextTerm(fakeDatasource, searchTerm)
			// if err == nil {
			// 	s.overwriteLastWordOnScreen(tabCompletion)
			// }
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

func (s *Shell) getLastWordFromBuffer() string {
	bufferString := s.buffer.String()
	if bufferString == "" {
		return ""
	}

	components := strings.Split(bufferString, " ")
	return components[len(components)-1]
}

func (s *Shell) getchar() ([]byte, error) {
	s.term.SetRaw()

	// not sure if this should be hardcoded as 256 chars
	// This limits pasting into the buffer to 256 chars
	bytes := make([]byte, 256)
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

	// delete everything in current buffers
	length := s.buffer.Len()
	for i := 0; i < length; i++ {
		s.term.Write(backspaceBytes)
	}

	s.flushBuffer()

	s.buffer.Write(bufferBytes)
	s.term.Write(bufferBytes)
}

func (s *Shell) overwriteLastWordOnScreen(word string) {
	bufferString := s.buffer.String()
	if bufferString == "" {
		s.buffer.Write([]byte(word))
		s.term.Write([]byte(word))
		return
	}

	components := strings.Split(bufferString, " ")
	lastWord := components[len(components)-1]

	// delete last word from buffers
	for i := 0; i < len(lastWord); i++ {
		s.term.Write(backspaceBytes)
	}

	length := len(bufferString) - len(lastWord)
	s.buffer.Truncate(length)

	s.buffer.Write([]byte(word))
	s.term.Write([]byte(word))
}

func exitHandler(args []string, s *Shell) (string, error) {
	s.Cleanup()
	os.Exit(0)
	return "", nil
}

// Cleanup does any work needed to cleanly close the shell
func (s *Shell) Cleanup() {
	s.term.Restore()
	s.term.Close()
}
