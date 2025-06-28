//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package term

import (
	"errors"
	"os"
	"syscall"
)

func GetTtySizeByName(fn string) (cols, rows int, err error)     { return getTtySize(fn) }
func GetTtySizeByFile(outf *os.File) (cols, rows int, err error) { return getDeviceSize(outf) }
func GetTtySizeByFd(fd uintptr) (cols, rows int, err error)      { return GetFdSize(fd) }

func errIsENOTTY(err error) bool {
	if errors.Is(err, syscall.ENOTTY) {
		return true
	}
	return false
}
