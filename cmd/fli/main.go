package main

import (
	"flag"
	"fmt"
	"os"

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
	s.AddCommand("ls", fli.lsHandler)
	s.AddCommand("pwd", fli.pwdHandler)
	s.AddCommand("cd", fli.cdHandler)

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

func (fli *Fli) lsHandler(args []string, s *shell.Shell) (string, error) {
	var p string

	if len(args) <= 1 {
		p = ""
	} else {
		p = args[1]
	}

	return fli.fStore.Ls(p)
}

func (fli *Fli) pwdHandler(args []string, s *shell.Shell) (string, error) {
	return fli.fStore.FirebaseURLFromWorkingDirectory("."), nil
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
