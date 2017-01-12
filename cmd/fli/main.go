package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sneakybueno/fli/fuego"
	"github.com/sneakybueno/fli/shell"
)

var fStore *fuego.FStore

func main() {
	// hijack Stdin
	stdin := shell.InitStdin()
	// XXX: doesn't always exit cleanly :/
	defer stdin.Cleanup()

	firebaseURL := "https://go-fli.firebaseio.com/"
	fStore = fuego.NewFStore(firebaseURL)

	fmt.Printf("Time to fli @ %s\n", fStore.WorkingDirectoryURL())
	fmt.Print(fStore.Prompt())

	for stdin.ReadNext() {
		input, ok := stdin.Text()
		if ok {
			m, err := processInput(input)
			if err != nil {
				fmt.Println(err)
			} else if m != "" {
				fmt.Println(m)
			}

			fmt.Print(fStore.Prompt())
			continue
		}

		// no input, check for keypresses we care about
		keyPress := stdin.KeyPress()
		switch keyPress {
		case shell.ArrowUp:
			fmt.Printf("arrow up")
		case shell.ArrowDown:
			fmt.Printf("arrow down")
		case shell.Tab:
			fmt.Printf("tab")
		}
	}
}

func processInput(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	components := strings.Split(input, " ")
	command := components[0]

	switch command {
	case "cd":
		var dir string
		if len(components) <= 1 {
			dir = ""
		} else {
			dir = components[1]
		}
		fStore.Cd(dir)
		return "", nil
	case "ls":
		return fStore.Ls()
	case "pwd":
		m := fmt.Sprintf("%s", fStore.WorkingDirectoryURL())
		return m, nil
	default:
		message := fmt.Sprintf("command not found: %s", command)
		return "", fliError(message)
	}
}

func fliError(message string) error {
	m := fmt.Sprintf("fli: %s", message)
	return errors.New(m)
}
