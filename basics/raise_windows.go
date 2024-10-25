package basics

import (
	"errors"
	"syscall"
)

// Raise an OS signal is not support on Windows.
func Raise(sig syscall.Signal) error {
	return errors.New("not supported")
}
