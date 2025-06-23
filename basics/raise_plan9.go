//go:build plan9
// +build plan9

package basics

import "github.com/hedzr/is/basics/syscall"

func Raise(sig syscall.Signal) error {
	return nil
}
