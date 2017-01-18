package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sneakybueno/fli/fuego"
	"github.com/sneakybueno/fli/shell"
)

func main() {
	s, err := shell.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// XXX: doesn't always exit cleanly :/
	defer s.Cleanup()

	var firebaseURL string
	var serviceAccountPath string

	flag.StringVar(&firebaseURL, "host", "", "Firebase database URL (Required)")
	flag.StringVar(&serviceAccountPath, "config", "", "Path to service account file (Required)")
	flag.Parse()

	if firebaseURL == "" || serviceAccountPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fStore, err := fuego.NewFStore(firebaseURL, serviceAccountPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Time to fli @ %s\n", fStore.WorkingDirectoryURL())
	fmt.Print(fStore.Prompt())

	for s.Next() {
		result, err := processInput(fStore, s.Input())
		if err != nil {
			fmt.Println(err)
		} else if result != "" {
			fmt.Println(result)
		}
		fmt.Print(fStore.Prompt())
	}

	if err = s.Error(); err != nil {
		fmt.Println(err)
	}
}

func processInput(fStore *fuego.FStore, input string) (string, error) {
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
		return fStore.WorkingDirectoryURL(), nil
	default:
		message := fmt.Sprintf("command not found: %s", command)
		return "", fliError(message)
	}
}

func fliError(message string) error {
	m := fmt.Sprintf("fli: %s", message)
	return errors.New(m)
}
