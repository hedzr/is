//go:build !windows && !plan9
// +build !windows,!plan9

package basics

import (
	"os"
	"syscall"
)

func Raise(sig syscall.Signal) error {
	return syscall.Kill(os.Getpid(), sig)
}
