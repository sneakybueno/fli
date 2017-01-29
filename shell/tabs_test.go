package shell_test

import (
	"testing"

	"github.com/sneakybueno/fli/shell"
	"github.com/stretchr/testify/assert"
)

func TestTabCompletionSearch(t *testing.T) {
	given := []string{"apple", "boom", "bueno", "test", "users"}

	result, err := shell.FindNextTerm(given, "b")
	assert.Nil(t, err)
	assert.Equal(t, result, "boom")

	result, err = shell.FindNextTerm(given, "bo")
	assert.Nil(t, err)
	assert.Equal(t, result, "boom")

	result, err = shell.FindNextTerm(given, "boom")
	assert.Nil(t, err)
	assert.Equal(t, result, "boom")

	result, err = shell.FindNextTerm(given, "bu")
	assert.Nil(t, err)
	assert.Equal(t, result, "bueno")

	result, err = shell.FindNextTerm(given, "boomShouldStillWork")
	assert.Nil(t, err)
	assert.Equal(t, result, "boom")

	result, err = shell.FindNextTerm(given, "p")
	assert.NotNil(t, err)
	assert.Equal(t, result, "")

	result, err = shell.FindNextTerm(given, " boom")
	assert.NotNil(t, err)
	assert.Equal(t, result, "")

	result, err = shell.FindNextTerm(given, "really long string that does not ex")
	assert.NotNil(t, err)
	assert.Equal(t, result, "")
}
