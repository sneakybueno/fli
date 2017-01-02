package fuego_test

import (
	"testing"

	"github.com/sneakybueno/fli/fuego"
)

const (
	firebaseTestingURL = "https://go-fli.firebaseio.com/"
)

func TestDirectoryCommands(t *testing.T) {
	fStore := fuego.NewFStore(firebaseTestingURL)

	expected := ""
	wd := fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}

	fStore.Cd("users")
	expected = "users"
	wd = fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}

	fStore.Cd("bueno/dev")
	expected = "users/bueno/dev"
	wd = fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}

	fStore.Cd("..")
	expected = "users/bueno"
	wd = fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}

	fStore.Cd("../corgi")
	expected = "users/corgi"
	wd = fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}

	fStore.Cd("../..")
	expected = ""
	wd = fStore.Wd()
	if wd != expected {
		t.Errorf("Expected %s, got %s", expected, wd)
	}
}
