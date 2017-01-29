package shell

import (
	"fmt"
	"sort"
	"strings"
)

type CommandHandler func(args []string, s *Shell) (string, error)

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

func (s *Shell) AddCommand(name string, handler CommandHandler) *Command {
	command := Command{
		Name:    name,
		Handler: handler,
	}

	s.commands = append(s.commands, command)
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
		fmt.Println(err)
		fmt.Print(s.prompt)

		return "", err
	}

	out, err := command.Handler(components, s)
	if err != nil {
		fmt.Println(err)
	} else if out != "" {
		fmt.Println(out)
	}

	fmt.Print(s.prompt)

	return out, err
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
