//go:build plan9
// +build plan9

package basics

import (
	"os"
)

func RaiseSyscallSignal(sig any) error {
	return raise(sig)
}

func raise(sig any) error {
	return nil
}

func raiseOsSig(sig os.Signal) error {
	return nil
}
