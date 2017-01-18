package shell_test

import (
	"testing"

	"github.com/sneakybueno/fli/shell"
)

func TestFindingCommands(t *testing.T) {
	s := &shell.Shell{}

	// Empty shell
	command, err := s.FindCommand("bad")
	if err == nil {
		t.Errorf("Expected an error and no command, got: %s", command.Name)
	}

	// Command length = 1
	s.AddCommand("one", nil)
	command, err = s.FindCommand("bad")
	if err == nil {
		t.Errorf("Expected an error and no command, got: %s", command.Name)
	}

	command, err = s.FindCommand("one")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}

	// Command length = 2
	s.AddCommand("two", nil)
	command, err = s.FindCommand("bad")
	if err == nil {
		t.Errorf("Expected an error and no command, got: %s", command.Name)
	}

	command, err = s.FindCommand("one")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}

	command, err = s.FindCommand("two")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}

	// Command length = 3
	s.AddCommand("three", nil)
	command, err = s.FindCommand("bad")
	if err == nil {
		t.Errorf("Expected an error and no command, got: %s", command.Name)
	}

	command, err = s.FindCommand("one")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}

	command, err = s.FindCommand("two")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}

	command, err = s.FindCommand("three")
	if err != nil {
		t.Errorf("Expected no error and command, got: %s", err)
	}
}
