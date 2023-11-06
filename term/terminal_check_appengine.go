//go:build appengine
// +build appengine

package term

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
