//go:build plan9 || appengine || wasm
// +build plan9 appengine wasm

package color

import (
	"os"
)

// Copyright © 2022 Atonal Authors
//

func (s *Cursor) CursorUp(n int) csiS      { return s.CSI('A', n) } // use color.Up() instead of this
func (s *Cursor) CursorDown(n int) csiS    { return s.CSI('B', n) } // use color.Down() instead of this
func (s *Cursor) CursorForward(n int) csiS { return s.CSI('C', n) } // use color.Right() instead of this
func (s *Cursor) CursorBack(n int) csiS    { return s.CSI('D', n) } // use color.Left() instead of this

func (s *Cursor) CursorNextLine(n int) csiS     { return s.CSI('E', n) }        // Moves cursor to beginning of the line n (default 1) lines down. (not ANSI.SYS)
func (s *Cursor) CursorPrevLine(n int) csiS     { return s.CSI('F', n) }        // Moves cursor to beginning of the line n (default 1) lines up. (not ANSI.SYS)
func (s *Cursor) CursorHorzCol(colAbs int) csiS { return s.CSI('G', colAbs) }   // Moves the cursor to column n (default 1).
func (s *Cursor) CursorPos(col, row int) csiS   { return s.CSI('H', col, row) } //
func (s *Cursor) CursorErase(n EraseTo) csiS    { return s.CSI('J', int(n)) }   // Erase in Display
func (s *Cursor) CursorEraseInLine(n int) csiS  { return s.CSI('K', n) }        // Erase in Line

func (s *Cursor) CursorScrollUp(n int) csiS   { return s.CSI('S', n) } // use color.ScrollUp() instead of this
func (s *Cursor) CursorScrollDown(n int) csiS { return s.CSI('T', n) } // use color.ScrollDown() instead of this

func (s *Cursor) CursorHorzVertPos(n, m int) csiS { return s.CSI('f', n, m) } // Horizontal Vertical Position
func (s *Cursor) CursorSGR(n int) csiS            { return s.CSI('m', n) }    // Select Graphic Rendition
func (s *Cursor) AUXPortOn() csiS                 { return s.CSI('i', 5) }    // AUX Port On
func (s *Cursor) AUXPortOff() csiS                { return s.CSI('i', 4) }    // AUX Port Off
func (s *Cursor) DSR() csiS                       { return s.CSI('n', 6) }    // Device Status Report

func (s *Cursor) CursorSavePos() csiS    { return s.CSI('s') } // Save Current Cursor Position
func (s *Cursor) CursorRestorePos() csiS { return s.CSI('u') } // Restore Current Cursor Position

//
//
//

func cursorUp(w Writer, n int) {
	writecsiseq(w, 'A', n)
	// var sb bytes.Buffer
	// var bb = []byte(aecHideCursor)
	// sb.Write(bb[0:2])
	// var ss = strconv.Itoa(n)
	// sb.WriteString(ss)
	// sb.WriteByte('A')
	// safeWrite(sb.Bytes())
	// // _, _ = fmt.Fprintf(Out, "%s[%dA", escape, n)
}

func cursorDown(w Writer, n int) {
	writecsiseq(w, 'B', n)
}

func cursorRight(w Writer, n int) {
	writecsiseq(w, 'C', n)
}

func cursorLeft(w Writer, n int) {
	writecsiseq(w, 'D', n)
	// var sb bytes.Buffer
	// var bb = []byte(aecHideCursor)
	// sb.Write(bb[0:2])
	// var ss = strconv.Itoa(n)
	// sb.WriteString(ss)
	// sb.WriteByte('D')
	// safeWrite(sb.Bytes())
	// // _, _ = fmt.Fprintf(Out, "%s[%dD", escape, n)
}

func cursorScrollUp(w Writer, n int)   { writecsiseq(w, 'S', n) }
func cursorScrollDown(w Writer, n int) { writecsiseq(w, 'T', n) }

func cursorSavePos(w Writer)    { writecsi(w, 's') }
func cursorRestorePos(w Writer) { writecsi(w, 'u') }

func cursorHorizontalAbsolute(w Writer, n int) {
	writecsiseq(w, 'G', n)
}

func hideCursor(w Writer) {
}

func showCursor(w Writer) {
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

func safeWrite(w Write, b []byte) (n int, e error) {
	return os.Stdout.Write(b)
}
