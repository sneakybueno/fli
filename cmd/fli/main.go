package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/sneakybueno/fli/fuego"
	"github.com/sneakybueno/fli/shell"
)

type Fli struct {
	fStore *fuego.FStore
}

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

	fmt.Printf("Time to fli @ %s\n", fStore.FirebaseURL)

	s, err := shell.Init(fStore.Prompt())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fli := &Fli{fStore: fStore}

	// Register command handlers
	s.AddCommand("hello", fli.helloHandler)

	s.AddCommand("cd", fli.cdHandler)
	s.AddCommand("find", fli.searchHandler)
	s.AddCommand("ls", fli.lsHandler)
	s.AddCommand("locate", fli.indexedSearchHandler)
	s.AddCommand("open", fli.openHandler)
	s.AddCommand("pwd", fli.pwdHandler)

	for s.Next() {
		s.Process(s.Input())
	}

	if err = s.Error(); err != nil {
		fmt.Println(err)
	}
}

func (fli *Fli) helloHandler(args []string, s *shell.Shell) (string, error) {
	return "Hello World -Fli", nil
}

func (fli *Fli) cdHandler(args []string, s *shell.Shell) (string, error) {
	var dir string

	if len(args) <= 1 {
		dir = ""
	} else {
		dir = args[1]
	}

	fli.fStore.Cd(dir)
	s.SetPrompt(fli.fStore.Prompt())

	return "", nil
}

func (fli *Fli) lsHandler(args []string, s *shell.Shell) (string, error) {
	var p string

	if len(args) <= 1 {
		p = ""
	} else {
		p = args[1]
	}

	return fli.fStore.Ls(p)
}

func (fli *Fli) openHandler(args []string, s *shell.Shell) (string, error) {
	var p string

	if len(args) <= 1 {
		p = ""
	} else {
		p = args[1]
	}

	url := fli.fStore.FirebaseURLFromWorkingDirectory(p)
	open.Start(url)

	message := fmt.Sprintf("opening (%s) in default browser", url)
	return message, nil
}

func (fli *Fli) pwdHandler(args []string, s *shell.Shell) (string, error) {
	return fli.fStore.FirebaseURLFromWorkingDirectory("."), nil
}

// Supports wild card searching
func (fli *Fli) indexedSearchHandler(args []string, s *shell.Shell) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("%s: [path] [key] [value]", args[0])
	}

	// Add support for searching from top level with ~
	// Allows user to seach in paths not based on cwd
	p := args[1]
	p = fli.fStore.BuildWorkingDirectoryPath(p)
	key := args[2]
	value := args[3]

	return fli.fStore.IndexedSearch(p, key, value)
}

func (fli *Fli) searchHandler(args []string, s *shell.Shell) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("%s: [path] [key] [value]", args[0])
	}

	// Add support for searching from top level with ~
	// Allows user to seach in paths not based on cwd
	p := args[1]
	p = fli.fStore.BuildWorkingDirectoryPath(p)
	key := args[2]
	value := args[3]

	return fli.fStore.Search(p, key, value)
}
