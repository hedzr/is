//go:build !aix && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !zos && !windows && !solaris && !plan9

package term

import (
	"errors"
	"os"
	"syscall"
)

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
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
