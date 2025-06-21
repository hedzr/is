package color

import (
	"fmt"
	"strconv"
)

// CSI code
type csiS struct {
	ch   uint8
	n, m int
	*Cursor
}

func (s csiS) Printf(format string, args ...any) *Cursor {
	_, _ = s.Cursor.sb.WriteString(csi)
	if s.n > 0 {
		_, _ = s.Cursor.sb.WriteString(strconv.Itoa(s.n))
	}
	if s.m > 0 {
		_ = s.Cursor.sb.WriteByte(';')
		_, _ = s.Cursor.sb.WriteString(strconv.Itoa(s.m))
	}
	_ = s.Cursor.sb.WriteByte(s.ch)

	return s.Cursor.Printf(format, args...)
}

func (s csiS) Echo(args ...string) *Cursor {
	s.Cursor.Echo(args...)
	return s.Cursor
}

func (s csiS) Println(args ...any) *Cursor {
	s.Cursor.Println(args...)
	return s.Cursor
}

func (s csiS) Print(args ...any) *Cursor {
	s.Cursor.Print(args...)
	return s.Cursor
}

func (s *Cursor) CSI(what uint8, n ...int) csiS {
	switch len(n) {
	case 0:
		return csiS{what, 0, 0, s}
	case 1:
		return csiS{what, n[0], 0, s}
	default:
		return csiS{what, n[0], n[1], s}
	}
}

// func (s *Cursor) CursorUp(n int) csiS      { return s.CSI('A', n) } // use color.Up() instead of this
// func (s *Cursor) CursorDown(n int) csiS    { return s.CSI('B', n) } // use color.Down() instead of this
// func (s *Cursor) CursorForward(n int) csiS { return s.CSI('C', n) } // use color.Right() instead of this
// func (s *Cursor) CursorBack(n int) csiS    { return s.CSI('D', n) } // use color.Left() instead of this

// func (s *Cursor) CursorNextLine(n int) csiS     { return s.CSI('E', n) }        // Moves cursor to beginning of the line n (default 1) lines down. (not ANSI.SYS)
// func (s *Cursor) CursorPrevLine(n int) csiS     { return s.CSI('F', n) }        // Moves cursor to beginning of the line n (default 1) lines up. (not ANSI.SYS)
// func (s *Cursor) CursorHorzCol(colAbs int) csiS { return s.CSI('G', colAbs) }   // Moves the cursor to column n (default 1).
// func (s *Cursor) CursorPos(col, row int) csiS   { return s.CSI('H', col, row) } //
// func (s *Cursor) CursorErase(n EraseTo) csiS    { return s.CSI('J', int(n)) }   // Erase in Display
// func (s *Cursor) CursorEraseInLine(n int) csiS  { return s.CSI('K', n) }        // Erase in Line

// func (s *Cursor) CursorScrollUp(n int) csiS   { return s.CSI('S', n) } // use color.ScrollUp() instead of this
// func (s *Cursor) CursorScrollDown(n int) csiS { return s.CSI('T', n) } // use color.ScrollDown() instead of this

// func (s *Cursor) CursorHorzVertPos(n, m int) csiS { return s.CSI('f', n, m) } // Horizontal Vertical Position
// func (s *Cursor) CursorSGR(n int) csiS            { return s.CSI('m', n) }    // Select Graphic Rendition
// func (s *Cursor) AUXPortOn() csiS                 { return s.CSI('i', 5) }    // AUX Port On
// func (s *Cursor) AUXPortOff() csiS                { return s.CSI('i', 4) }    // AUX Port Off
// func (s *Cursor) DSR() csiS                       { return s.CSI('n', 6) }    // Device Status Report

// func (s *Cursor) CursorSavePos() csiS    { return s.CSI('s') } // Save Current Cursor Position
// func (s *Cursor) CursorRestorePos() csiS { return s.CSI('u') } // Restore Current Cursor Position

//
// Special Cursor Operations
//

func (s *Cursor) Flush() *Cursor {
	if s.sb.Len() > 0 {
		fmt.Fprint(s.w, s.String())
		s.Reset()
	}
	return s
}

func (s *Cursor) HorizontalAbsolute(n int) *Cursor {
	s.Flush()
	cursorHorizontalAbsolute(s.w, n)
	return s
}

// =color.Up()
func (s *Cursor) Up(n int) *Cursor {
	s.Flush()
	cursorUp(s.w, n)
	return s
}

// =color.Down()
func (s *Cursor) Down(n int) *Cursor {
	s.Flush()
	cursorDown(s.w, n)
	return s
}

// =color.Right()
func (s *Cursor) Right(n int) *Cursor {
	s.Flush()
	cursorRight(s.w, n)
	return s
}

// =color.Left()
func (s *Cursor) Left(n int) *Cursor {
	s.Flush()
	cursorLeft(s.w, n)
	return s
}

// =color.ScrollUp()
func (s *Cursor) ScrollUp(n int) *Cursor {
	s.Flush()
	cursorScrollUp(s.w, n)
	return s
}

// =color.ScrollDown()
func (s *Cursor) ScrollDown(n int) *Cursor {
	s.Flush()
	cursorScrollDown(s.w, n)
	return s
}

// =color.SavePos()
func (s *Cursor) SavePos() *Cursor {
	s.Flush()
	cursorSavePos(s.w)
	return s
}

// =color.RestorePos()
func (s *Cursor) RestorePos() *Cursor {
	s.Flush()
	cursorRestorePos(s.w)
	return s
}
