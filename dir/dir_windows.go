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
	ts := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	time.Unix(0, ts.CreationTime.Nanoseconds())
	return
}

// FileAccessedTime return the creation time of a file
func FileAccessedTime(fileInfo os.FileInfo) (tm time.Time) {
	ts := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	time.Unix(0, ts.LastAccessTime.Nanoseconds())
	return
}

// FileModifiedTime return the creation time of a file
func FileModifiedTime(fileInfo os.FileInfo) (tm time.Time) {
	ts := fileInfo.Sys().(*syscall.Win32FileAttributeData)
	time.Unix(0, ts.LastWriteTime.Nanoseconds())
	return
}

func timeSpecToTime(ts syscall.Timespec) time.Time {
	// TODO ts.Sec is not ok for plan9
	return time.Unix(ts.Sec, ts.Nsec)
}
