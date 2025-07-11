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

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	cols, rows, _ = getFdSize(os.Stdout.Fd())
	// if cols == 0 && rows == 0 {
	// 	cols, rows, _ = getFdSize(uintptr(syscall.Stdout))
	// 	if cols == 0 && rows == 0 {
	// 		// st, _ := GetState(os.Stdout.Fd())
	// 		ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	// 		if err != nil {
	// 			return 0, 0
	// 		}
	// 		return int(ws.Col), int(ws.Row)
	// 	}
	// }
	// // cols, rows, _ = getTtySize("/dev/tty")
	return
}

func getTtySize(device string) (cols, rows int, err error) {
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

	var outf *os.File
	outf, err = os.OpenFile(device, os.O_WRONLY, 0)
	if err != nil {
		slog.Warn("CANNOT getTtySize", "err", err)
		err = nil
		return
	}
	defer outf.Close()
	return getDeviceSize(outf)
}

func getDeviceSize(outf *os.File) (cols, rows int, err error) {
	out := outf.Fd()
	return getFdSize(out)
}

func getFdSize(fd uintptr) (cols, rows int, err error) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}
	res, _, errno := syscall.Syscall(syscall.SYS_IOCTL, //nolint:dogsled //like it
		uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	if int(res) == -1 {
		err = fmt.Errorf("[getTtySize] cannot get terminal size. res = %v, err = %v", res, errno.Error())
		// slog.Error("[GetTtySize] cannot get terminal size", "err", err)
	}
	cols, rows = int(sz.cols), int(sz.rows)
	return
}

func isDoubleClickRun() bool { return false }
