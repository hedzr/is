//go:build !windows
// +build !windows

package basics

import (
	"os"
	"syscall"
)

func Raise(sig syscall.Signal) error {
	return syscall.Kill(os.Getpid(), sig)
}
