//go:build !windows && !plan9
// +build !windows,!plan9

package basics

import (
	"os"
	"syscall"
)

func RaiseSyscallSignal(sig syscall.Signal) error {
	return raise(sig)
}

func raise(sig syscall.Signal) error {
	return syscall.Kill(os.Getpid(), sig)
}

func raiseOsSig(sig os.Signal) error {
	if sigx, ok := sig.(syscall.Signal); ok {
		return syscall.Kill(os.Getpid(), sigx)
	}
	return nil
}
