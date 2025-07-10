//go:build windows
// +build windows

package basics

import (
	"errors"
	"os"
	"syscall"
)

func RaiseSyscallSignal(sig syscall.Signal) error {
	return raise(sig)
}

// Raise an OS signal is not support on Windows.
func raise(sig syscall.Signal) error {
	_ = sig
	return errors.New("not supported")
}

func raiseOsSig(sig os.Signal) error {
	_ = sig
	_ = os.Kill
	return errors.New("not supported")
}

const SIG_USR1 = syscall.SIGALRM
