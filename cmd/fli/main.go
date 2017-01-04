package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"flag"

	"github.com/sneakybueno/fli/fuego"
)

func main() {
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		m, err := processInput(fStore, input)
		if err != nil {
			fmt.Println(err)
		} else if m != "" {
			fmt.Println(m)
		}

		fmt.Print(fStore.Prompt())
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
