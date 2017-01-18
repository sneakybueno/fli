package shell

import (
	"bytes"
	"strings"
)

func nonEmpty(input string) bool {
	// TODO: figure out how to properly check for strings that are just whitespace
	if strings.TrimSpace(input) == "" {
		return false
	}
	return true
}

// ASCII Matching
// ----------------------------------------------------------------------------

func isEnter(b []byte) bool {
	return bytes.Equal(b, []byte{13})
}

func isDelete(b []byte) bool {
	return bytes.Equal(b, []byte{127})
}

func isTab(b []byte) bool {
	return bytes.Equal(b, []byte{9})
}

func isEsc(b []byte) bool {
	return bytes.Equal(b, []byte{27})
}

func isArrowUp(b []byte) bool {
	return bytes.Equal(b, []byte{27, 91, 65})
}

func isArrowRight(b []byte) bool {
	return bytes.Equal(b, []byte{27, 91, 67})
}

func isArrowDown(b []byte) bool {
	return bytes.Equal(b, []byte{27, 91, 66})
}

func isArrowLeft(b []byte) bool {
	return bytes.Equal(b, []byte{27, 91, 68})
}

func isCtrlC(b []byte) bool {
	return bytes.Equal(b, []byte{3})
}
