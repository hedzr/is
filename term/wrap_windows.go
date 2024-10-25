//go:build windows && !nacl
// +build windows,!nacl

// Copyright © 2020 Hedzr Yeh.

package term

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/crypto/ssh/terminal"
)

// // ReadPassword reads the password from stdin with safe protection
// func ReadPassword() (text string, err error) {
// 	var bytePassword []byte
// 	if bytePassword, err = terminal.ReadPassword(0); err == nil {
// 		fmt.Println() // it's necessary to add a new line after user's input
// 		text = string(bytePassword)
// 	} else {
// 		fmt.Println() // it's necessary to add a new line after user's input
// 	}
// 	return
// }

func GetTtySizeByName(fn string) (cols, rows int, err error) {
	cols, rows = GetTtySize()
	return
}

func GetTtySizeByFile(outf *os.File) (cols, rows int, err error) { return getDeviceSize(outf) }
func GetTtySizeByFd(fd uintptr) (cols, rows int, err error)      { return GetFdSize(fd) }

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	// return 0, 0
	cols, rows, _ = terminal.GetSize(0) // https://stackoverflow.com/a/45422726/6375060
	return
}

func getDeviceSize(outf *os.File) (cols, rows int, err error) {
	out := outf.Fd()
	return getFdSize(out)
}

func getFdSize(fd uintptr) (cols, rows int, err error) {
	cols, rows, _ = terminal.GetSize(int(fd))
	return
}

// isDoubleClickRun detects
// if windows golang executable file is running via double click or from cmd/shell terminator
//
// - https://stackoverflow.com/questions/8610489/distinguish-if-program-runs-by-clicking-on-the-icon-typing-its-name-in-the-cons?rq=1
// - https://github.com/shirou/w32/blob/master/kernel32.go
// - https://github.com/kbinani/win/blob/master/kernel32.go#L3268
// - win.GetConsoleProcessList(new(uint32), win.DWORD(2))
//
// From: https://gist.github.com/yougg/213250cc04a52e2b853590b06f49d865
//
// Sample code:
//
//	func main() {
//		clickRun := isDoubleClickRun()
//		fmt.Println("Double click run:", clickRun)
//		if clickRun {
//			fmt.Print("press Enter to exit")
//			var b byte
//			_, _ = fmt.Scanf("%v", &b)
//		}
//	}
func isDoubleClickRun() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}
