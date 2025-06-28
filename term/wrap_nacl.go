//go:build nacl
// +build nacl

// Copyright © 2020 Hedzr Yeh.

package term

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

// // ReadPassword reads the password from stdin with safe protection
// func ReadPassword() (text string, err error) {
// 	return stringtool.RandomStringPure(9), nil
// }

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	var sz struct {
		rows, cols, xPixels, yPixels uint16
	}
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&sz)))
	cols, rows = int(sz.cols), int(sz.rows)
	// cols, rows = 32768, 43
	return
}

func isDoubleClickRun() bool { return false }

func GetTtySizeByName(fn string) (cols, rows int, err error)     { return getTtySize(fn) }
func GetTtySizeByFile(outf *os.File) (cols, rows int, err error) { return getDeviceSize(outf) }
func GetTtySizeByFd(fd uintptr) (cols, rows int, err error)      { return GetFdSize(fd) }

func getTtySize(fn string) (cols, rows int, err error)        { return }
func getDeviceSize(outf *os.File) (cols, rows int, err error) { return }
func getFdSize(fd uintptr) (cols, rows int, err error)        { return }

func errIsENOTTY(err error) bool {
	if errors.Is(err, syscall.ENOTTY) {
		return true
	}
	return false
}
