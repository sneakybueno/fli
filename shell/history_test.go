package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoryEmpty(t *testing.T) {
	h := InitCmdHistory(5)
	assert.Empty(t, h.Next())
	assert.Empty(t, h.Prev())
	assert.Empty(t, h.All())
}

func TestHistoryCycle(t *testing.T) {
	h := InitCmdHistory(5)

	// add some cmds
	cmds := []string{"ls", "pwd", "cd .."}
	for _, cmd := range cmds {
		h.Add(cmd)
	}
	assert.Len(t, h.All(), len(cmds))

	// check cycling backwards
	for i := len(cmds) - 1; i >= 0; i-- {
		assert.Equal(t, cmds[i], h.Prev())
	}
	// should go back to the end of the array
	assert.Equal(t, cmds[len(cmds)-1], h.Prev())

	// check cycling forwards
	// note: h.All() resets the history offset
	t.Logf("History: %+v", h.All())
	for i := 0; i < len(cmds)*2; i++ {
		assert.Equal(t, cmds[i%len(cmds)], h.Next())
	}
}

func TestHistoryAdd(t *testing.T) {
	h := InitCmdHistory(5)
	h.Add("ls")
	cmds := []string{"ls"}

	// Add()ing anywhere in the cycle shouldn't affect where the cmd is added
	// since Add() resets the history offset
	h.Next()
	h.Add("cd bueno")
	cmds = append(cmds, "cd bueno")

	h.Next()
	h.Next()
	h.Add("cd utils")
	cmds = append(cmds, "cd utils")

	for i := len(cmds) - 1; i >= 0; i-- {
		assert.Equal(t, cmds[i], h.Prev())
	}
}

func TestHistoryOverflow(t *testing.T) {
	h := InitCmdHistory(5)
	cmds := []string{"ls", "pwd", "cd..", "cd bueno", "cd utils"}
	for _, cmd := range cmds {
		h.Add(cmd)
	}
	// all good, we only added 5
	assert.Len(t, h.All(), len(cmds))

	// add a 6th cmd
	h.Add("cd etc")
	assert.Len(t, h.All(), h.max)

	// check that the first cmd was bumped out
	assert.Equal(t, "cd etc", h.Prev())
	for i := len(cmds) - 1; i >= 1; i-- {
		assert.Equal(t, cmds[i], h.Prev())
	}
	// back the beginning
	assert.Equal(t, "cd etc", h.Prev())
}
