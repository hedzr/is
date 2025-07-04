/*
 * Copyright © 2021 Hedzr Yeh.
 */

package dir

import (
	"os"
	"syscall"
	"time"
)

// FileCreatedTime return the creation time of a file
func FileCreatedTime(fileInfo os.FileInfo) (tm time.Time) {
	if ts, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		tm = timeSpecToTime(ts.Ctimespec)
	}
	return
}

// FileAccessedTime return the creation time of a file
func FileAccessedTime(fileInfo os.FileInfo) (tm time.Time) {
	if ts, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		tm = timeSpecToTime(ts.Atimespec)
	}
	return
}

// FileModifiedTime return the creation time of a file
func FileModifiedTime(fileInfo os.FileInfo) (tm time.Time) {
	if ts, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		tm = timeSpecToTime(ts.Mtimespec)
	}
	return
}

func timeSpecToTime(ts syscall.Timespec) time.Time {
	// TODO ts.Sec is not ok for plan9
	return time.Unix(ts.Sec, ts.Nsec)
}
