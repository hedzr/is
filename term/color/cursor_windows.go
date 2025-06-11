// Copyright © 2022 Atonal Authors
//

//go:build windows
// +build windows

package color

import (
	"io"
	"syscall"
	"unsafe"
)

type (
	SHORT int16
	WORD  uint16

	SMALL_RECT struct {
		Left   SHORT
		Top    SHORT
		Right  SHORT
		Bottom SHORT
	}

	COORD struct {
		X SHORT
		Y SHORT
	}

	CONSOLE_SCREEN_BUFFER_INFO struct {
		Size              COORD
		CursorPosition    COORD
		Attributes        WORD
		Window            SMALL_RECT
		MaximumWindowSize COORD
	}
	CONSOLE_CURSOR_INFO struct {
		Size    uint32
		Visible int32
	}
)

var (
	getConsoleScreenBufferInfoProc *syscall.LazyProc
	getConsoleCursorPositionProc   *syscall.LazyProc
	setConsoleCursorPositionProc   *syscall.LazyProc
	getConsoleCursorInfoProc       *syscall.LazyProc
	setConsoleCursorInfoProc       *syscall.LazyProc
)

func init() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleCursorInfoProc = kernel32.NewProc("GetConsoleCursorInfo")
	setConsoleCursorInfoProc = kernel32.NewProc("SetConsoleCursorInfo")
	getConsoleScreenBufferInfoProc = kernel32.NewProc("GetConsoleScreenBufferInfo")
	getConsoleCursorPositionProc = kernel32.NewProc("GetConsoleCursorPosition")
	setConsoleCursorPositionProc = kernel32.NewProc("SetConsoleCursorPosition")
}

// checkError evaluates the results of a Windows API call and returns the error if it failed.
func checkError(r1, r2 uintptr, err error) error {
	// Windows APIs return non-zero to indicate success
	if r1 != 0 {
		return nil
	}

	// Return the error if provided, otherwise default to EINVAL
	if err != nil {
		return err
	}
	return syscall.EINVAL
}

// coordToPointer converts a COORD into a uintptr (by fooling the type system).
func coordToPointer(c COORD) uintptr {
	// Note: This code assumes the two SHORTs are correctly laid out; the "cast" to uint32 is just to get a pointer to pass.
	return uintptr(*((*uint32)(unsafe.Pointer(&c))))
}

func getStdHandle(stdhandle int) (uintptr, error) {
	handle, err := syscall.GetStdHandle(stdhandle)
	if err != nil {
		return 0, err
	}
	return uintptr(handle), nil
}

// GetConsoleScreenBufferInfo retrieves information about the specified console screen buffer.
// See http://msdn.microsoft.com/en-us/library/windows/desktop/ms683171(v=vs.85).aspx.
func getConsoleScreenBufferInfo(handle uintptr) (info *CONSOLE_SCREEN_BUFFER_INFO, err error) {
	info = &CONSOLE_SCREEN_BUFFER_INFO{}
	err = checkError(getConsoleScreenBufferInfoProc.Call(handle, uintptr(unsafe.Pointer(info)), 0))
	return
}

// SetConsoleCursorPosition location of the console cursor.
// See https://msdn.microsoft.com/en-us/library/windows/desktop/ms686025(v=vs.85).aspx.
func setConsoleCursorPosition(handle uintptr, coord COORD) error {
	r1, r2, err := setConsoleCursorPositionProc.Call(handle, coordToPointer(coord))
	// use(coord)
	return checkError(r1, r2, err)
}

func getConsoleCursorPosition(handle uintptr) (coord COORD, err error) {
	err = checkError(getConsoleCursorPositionProc.Call(handle, coordToPointer(coord)))
	return
}

func showHideCursor(w Writer, visible bool) (err error) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// var handle uintptr
	// handle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE) // syscall.Handle(os.Stdout.Fd())
	// if err != nil {
	// 	return
	// }

	var cursorInfo CONSOLE_CURSOR_INFO
	err = checkError(getConsoleCursorInfoProc.Call(stdoutHandle, uintptr(unsafe.Pointer(&cursorInfo))))
	if err != nil {
		return
	}

	cursorInfo.Visible = func() int32 {
		if visible {
			return 1
		} else {
			return 0
		}
	}()

	err = checkError(setConsoleCursorInfoProc.Call(stdoutHandle, uintptr(unsafe.Pointer(&cursorInfo))))
	return
}

func hideCursor(w Writer) error {
	return showHideCursor(w, false)
}

func showCursor(w Writer) error {
	return showHideCursor(w, true)
}

func cursorUp(w Writer, n int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// var err error
	// var stdoutHandle uintptr
	// stdoutHandle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	// if err != nil {
	// 	return
	// }

	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	y := consoleInfo.CursorPosition.Y - SHORT(n)
	setConsoleCursorPosition(stdoutHandle, COORD{X: consoleInfo.CursorPosition.X, Y: y})
}

func cursorDown(w Writer, n int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// var err error
	// var stdoutHandle uintptr
	// stdoutHandle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	// if err != nil {
	// 	return
	// }

	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	y := consoleInfo.CursorPosition.Y + SHORT(n)
	setConsoleCursorPosition(stdoutHandle, COORD{X: consoleInfo.CursorPosition.X, Y: y})
}

func cursorRight(w Writer, n int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// var err error
	// var stdoutHandle uintptr
	// stdoutHandle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	// if err != nil {
	// 	return
	// }

	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	x := consoleInfo.CursorPosition.X + SHORT(n)
	setConsoleCursorPosition(stdoutHandle, COORD{X: x, Y: consoleInfo.CursorPosition.Y})
}

func cursorLeft(w Writer, n int) {
	// var err error
	// var stdoutHandle uintptr
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// stdoutHandle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	// if err != nil {
	// 	return
	// }

	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	x := consoleInfo.CursorPosition.X - SHORT(n)
	setConsoleCursorPosition(stdoutHandle, COORD{X: x, Y: consoleInfo.CursorPosition.Y})
}

func cursorScrollUp(w Writer, n int)   { writecsiseq(w, 'S', n) }
func cursorScrollDown(w Writer, n int) { writecsiseq(w, 'T', n) }

func cursorSavePos(w Writer)    { writecsi(w, 's') }
func cursorRestorePos(w Writer) { writecsi(w, 'u') }

func cursorHorizontalAbsolute(w Writer, n int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	// stdoutHandle, err = getStdHandle(syscall.STD_OUTPUT_HANDLE)

	// var csbi consoleScreenBufferInfo
	// _, _, _ = procGetConsoleScreenBufferInfo.Call(stdoutHandle, uintptr(unsafe.Pointer(&csbi)))

	csbi, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	var cursor COORD
	cursor.X = SHORT(n)
	cursor.Y = csbi.CursorPosition.Y

	if csbi.Size.X < cursor.X {
		cursor.X = csbi.Size.X
	}

	_ = setConsoleCursorPosition(stdoutHandle, cursor)
}

func writecsi(out Writer, what rune) {
	// var bb = []byte(aecHideCursor)
	// _, _ = out.Write(bb[0:2])
	// _, _ = out.Write([]byte{byte(what)})
	_, _ = out, what
}

func writecsiseq(out Writer, what rune, n int) {
	// var sb bytes.Buffer
	// var bb = []byte(aecHideCursor)
	// _, _ = sb.Write(bb[0:2])
	// var ss = strconv.Itoa(n)
	// _, _ = sb.WriteString(ss)
	// _ = sb.WriteByte(byte(what))
	// _, _ = out.Write(sb.Bytes())
}

func safeWrite(w io.Writer, b []byte) (n int, e error) {
	return w.Write(b)
}
