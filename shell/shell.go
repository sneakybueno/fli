package shell

import (
	"fmt"
	"strings"
)

// Shell struct for keeping track of shell things
type Shell struct {
	stdin   *Stdin
	history *CmdHistory
	input   string
}

// Init creates a shell-like env
func Init() *Shell {
	return &Shell{
		stdin:   InitStdin(),
		history: InitCmdHistory(50),
	}
}

// Next returns true if there is more input
func (s *Shell) Next() bool {
	for s.stdin.ReadNext() {
		switch s.stdin.KeyPress() {
		case Enter:
			s.input = s.stdin.Flush()
			if nonEmpty(s.input) {
				s.history.Add(s.input)
			}
			return true
		case ArrowUp:
			s.stdin.Override(s.history.Prev())
			fakeRedraw(s.stdin.Peek())
		case ArrowDown:
			s.stdin.Override(s.history.Next())
			fakeRedraw(s.stdin.Peek())
		case Delete:
			fakeRedraw(s.stdin.Peek())
		case Tab:
			fmt.Printf("[tab]")
		}
	}

	return false
}

// Input returns any available input
func (s *Shell) Input() string {
	return s.input
}

// Cleanup does any work needed to cleanly close the shell
func (s *Shell) Cleanup() {
	s.stdin.Restore()
}

// because how the fuck do you do this?
func fakeRedraw(input string) {
	if nonEmpty(input) {
		fmt.Printf(" -> %s", input)
	}
}

func nonEmpty(input string) bool {
	// TODO: figure out how to properly check for strings that are just whitespace
	if strings.TrimSpace(input) == "" {
		return false
	}
	return true
}
