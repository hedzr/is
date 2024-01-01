//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package term

import (
	"os"
)

func GetTtySizeByName(fn string) (cols, rows int, err error)     { return getTtySize(fn) }
func GetTtySizeByFile(outf *os.File) (cols, rows int, err error) { return getDeviceSize(outf) }
func GetTtySizeByFd(fd uintptr) (cols, rows int, err error)      { return GetFdSize(fd) }
