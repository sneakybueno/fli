package shell

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"sort"

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

	input string
	err   error
}

type CommandHandler func(args []string) (string, error)

type Command struct {
	Name    string
	Handler CommandHandler
}

type Commands []Command

func (commands Commands) Len() int {
	return len(commands)
}

func (commands Commands) Less(i, j int) bool {
	return commands[i].Name < commands[j].Name
}

func (commands Commands) Swap(i, j int) {
	commands[i], commands[j] = commands[j], commands[i]
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

func (s *Shell) AddCommand(name string, handler CommandHandler) *Command {
	command := Command{
		Name:    name,
		Handler: handler,
	}

	s.commands = append(s.commands, command)
	// need to sort here
	sort.Sort(s.commands)

	return &command
}

func (s *Shell) Process(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	components := strings.Split(input, " ")
	commandString := components[0]

	command, err := s.FindCommand(commandString)
	if err != nil {
		return "", err
	}

	out, err := command.Handler(components)
	if err != nil {
		return "", err
	}

	return out, nil
}

// Binary search for command based on name
func (s *Shell) FindCommand(commandName string) (*Command, error) {
	l := 0
	r := len(s.commands) - 1

	for l <= r {
		mid := (r + l) / 2
		c := s.commands[mid]

		comparison := strings.Compare(commandName, c.Name)
		if comparison == 0 {
			return &c, nil
		}

		if comparison < 0 {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	return nil, fmt.Errorf("shell: command not found: %s", commandName)
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
