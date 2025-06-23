// Copyright © 2022 Atonal Authors
//

//go:build windows
// +build windows

package color

import (
	"context"
	"io"
	"log/slog"
	"os"
	"syscall"
	"unsafe"
)

func (s *Cursor) pCSI(suffix byte, args ...int) csiS {
	s.Flush()
	switch suffix {
	case 'A':
		cursorUp(s.w, args[0])
	case 'B':
		cursorDown(s.w, args[0])
	case 'C':
		cursorRight(s.w, args[0])
	case 'D':
		cursorLeft(s.w, args[0])
	case 'E', 'F', 'J', 'K':
		cursorUp(s.w, args[0])
	case 'G':
		cursorHorizontalAbsolute(s.w, args[0])
	case 'H':
		cursorSetPos(s.w, args[0], args[1])
	case 'S':
		cursorScrollUp(s.w, args[0])
	case 'T':
		cursorScrollDown(s.w, args[0])
	case 'f': //horz vert pos
	case 's':
		cursorSavePos(s.w)
	case 'u':
		cursorRestorePos(s.w)
	}

	return s.CSI(suffix, args...)
}

func (s *Cursor) Up(n int) csiS      { return s.pCSI('A', n) } // use color.Up() instead of this
func (s *Cursor) Down(n int) csiS    { return s.pCSI('B', n) } // use color.Down() instead of this
func (s *Cursor) Forward(n int) csiS { return s.pCSI('C', n) } // use color.Right() instead of this
func (s *Cursor) Back(n int) csiS    { return s.pCSI('D', n) } // use color.Left() instead of this

func (s *Cursor) NextLine(n int) csiS      { return s.pCSI('E', n) }        // Moves cursor to beginning of the line n (default 1) lines down. (not ANSI.SYS)
func (s *Cursor) PrevLine(n int) csiS      { return s.pCSI('F', n) }        // Moves cursor to beginning of the line n (default 1) lines up. (not ANSI.SYS)
func (s *Cursor) HorzCol(colAbs int) csiS  { return s.pCSI('G', colAbs) }   // Moves the cursor to column n (default 1).
func (s *Cursor) MoveTo(col, row int) csiS { return s.pCSI('H', col, row) } //
func (s *Cursor) Erase(n EraseTo) csiS     { return s.pCSI('J', int(n)) }   // Erase in Display
func (s *Cursor) EraseInLine(n int) csiS   { return s.pCSI('K', n) }        // Erase in Line

func (s *Cursor) ScrollUp(n int) csiS   { return s.pCSI('S', n) } // use color.ScrollUp() instead of this
func (s *Cursor) ScrollDown(n int) csiS { return s.pCSI('T', n) } // use color.ScrollDown() instead of this

func (s *Cursor) HorzVertPos(n, m int) csiS { return s.pCSI('f', n, m) } // Horizontal Vertical Position

func (s *Cursor) SGR(n int) csiS   { return s.pCSI('m', n) } // Select Graphic Rendition
func (s *Cursor) AUXPortOn() csiS  { return s.pCSI('i', 5) } // AUX Port On
func (s *Cursor) AUXPortOff() csiS { return s.pCSI('i', 4) } // AUX Port Off
func (s *Cursor) DSR() csiS        { return s.pCSI('n', 6) } // Device Status Report

func (s *Cursor) SavePos() csiS    { return s.pCSI('s') } // Save Current Cursor Position
func (s *Cursor) RestorePos() csiS { return s.pCSI('u') } // Restore Current Cursor Position

func (s *Cursor) CursorGet(ctx context.Context, pos *CursorPos) *Cursor {
	if pos != nil {
		var stdoutHandle uintptr = uintptr(syscall.Handle(os.Stdout.Fd()))
		consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
		if err != nil {
			slog.Error("CursorGet() [windows] getConsoleScreenBufferInfo failed", "err", err)
			return s
		}
		pos.Row = int(consoleInfo.CursorPosition.Y)
		pos.Col = int(consoleInfo.CursorPosition.X)
	}
	return s
}

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
	if x < 1 {
		x = 1
	}
	setConsoleCursorPosition(stdoutHandle, COORD{X: x, Y: consoleInfo.CursorPosition.Y})
}

func cursorScrollUp(w Writer, n int) {
	// writecsiseq(w, 'S', n)
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	y := consoleInfo.CursorPosition.Y - SHORT(n)
	if y < 1 {
		y = 1
	}
	setConsoleCursorPosition(stdoutHandle, COORD{X: consoleInfo.CursorPosition.X, Y: y})
}

func cursorScrollDown(w Writer, n int) {
	// writecsiseq(w, 'T', n)
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
	consoleInfo, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	y := consoleInfo.CursorPosition.Y + SHORT(n)
	setConsoleCursorPosition(stdoutHandle, COORD{X: consoleInfo.CursorPosition.X, Y: y})
}

var savedCursorPos []COORD

func cursorSavePos(w Writer) {
	// writecsi(w, 's')
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))

	csbi, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}
	coord := csbi.Size
	savedCursorPos = append(savedCursorPos, coord)
}

func cursorRestorePos(w Writer) {
	// writecsi(w, 'u')
	if l := len(savedCursorPos); l > 0 {
		coord := savedCursorPos[l-1]
		savedCursorPos = savedCursorPos[0 : l-1]
		var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))
		_ = setConsoleCursorPosition(stdoutHandle, coord)
	}
}

func cursorGetPos(w Writer) (row, col int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))

	csbi, err := getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		return
	}

	col, row = int(csbi.Size.X), int(csbi.Size.Y)
	return
}

func cursorSetPos(w Writer, row, col int) {
	var stdoutHandle uintptr = uintptr(syscall.Handle(w.Fd()))

	var cursor COORD
	cursor.X = SHORT(col)
	cursor.Y = SHORT(row)

	_ = setConsoleCursorPosition(stdoutHandle, cursor)
}

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
