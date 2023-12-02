//go:build darwin
// +build darwin

// Copyright © 2020 Hedzr Yeh.

package term

import (
	"fmt"
	"log/slog"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/crypto/ssh/terminal"
)

// NOTE:
//   SA1019: package golang.org/x/crypto/ssh/terminal is deprecated: this package moved to golang.org/x/term.
// Here we keep old reference for backward-compatibility to go1.11 (even lower)

//

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	var bytePassword []byte
	if bytePassword, err = terminal.ReadPassword(syscall.Stdin); err == nil {
		fmt.Println() // it's necessary to add a new line after user's input
		text = string(bytePassword)
	} else {
		fmt.Println() // it's necessary to add a new line after user's input
	}
	return
}

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}

	// var err error
	//
	// if runtime.GOOS == "openbsd" || runtime.GOOS == "freebsd" {
	// 	out, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	in = int(out.Fd())
	// } else {
	// 	out, err = os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	in, err = syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// out := os.Stdout.Fd()

	var err error
	var outf *os.File
	outf, err = os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		slog.Error("err", "err", err)
		return
	}
	defer outf.Close()
	out := outf.Fd()

	res, _, errno := syscall.Syscall(syscall.SYS_IOCTL, //nolint:dogsled //like it
		uintptr(out),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	if int(res) == -1 {
		slog.Error("[GetTtySize] cannot get terminal size", "err", errno.Error(), "res", res)
	}
	cols, rows = int(sz.cols), int(sz.rows)
	return
}
