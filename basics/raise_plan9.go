//go:build plan9
// +build plan9

package basics

import (
	"os"
	"syscall"
)

func RaiseSyscallSignal(sig any) error {
	return raise(sig)
}

func raise(sig any) error {
	return nil
}

func raiseOsSig(sig os.Signal) error {
	_ = os.Kill
	return nil
}

const SIG_USR1 = syscall.Note("usr1")
