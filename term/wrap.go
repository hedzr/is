// go: // build (dragonfly || freebsd || linux || netbsd || openbsd || aix || arm_linux || solaris) && !nacl && !plan9
// +b // uild dragonfly freebsd linux netbsd openbsd aix arm_linux solaris
// +b // uild !nacl
// +b // uild !plan9

// Copyright © 2020 Hedzr Yeh.

package term

import (
	"os"
)

// NOTE:
//   SA1019: package golang.org/x/crypto/ssh/terminal is deprecated: this package moved to golang.org/x/term.
// Here we keep old reference for backward-compatibility to go1.11 (even lower)

//

// // ReadPassword reads the password from stdin with safe protection
// func ReadPassword() (text string, err error) {
// 	var bytePassword []byte
// 	if bytePassword, err = term.ReadPassword(syscall.Stdin); err == nil {
// 		fmt.Println() // it's necessary to add a new line after user's input
// 		text = string(bytePassword)
// 	} else {
// 		fmt.Println() // it's necessary to add a new line after user's input
// 	}
// 	return
// }

// IsStartupByDoubleClick detects
// if windows golang executable file is running via double click or from cmd/shell terminator
func IsStartupByDoubleClick() bool {
	return isDoubleClickRun()
}

// IsCharDevice detect a file if it's a character device, for unix.
// Specially used for test if terminal under darwin, since new macOS has a different os.Stdout.
func IsCharDevice(f *os.File) bool {
	stat, _ := f.Stat()
	return (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

func GetFdSize(fd uintptr) (cols, rows int, err error) {
	return getFdSize(fd)
}
