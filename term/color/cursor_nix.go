//go:build !windows && !plan9 && !appengine && !wasm
// +build !windows,!plan9,!appengine,!wasm

package color

// Copyright © 2022 Atonal Authors

import (
	"bytes"
	"strconv"
)

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


// CursorGet try to get the current cursor position via ansi escaped sequences.
//
// Now it works ok for unix, linux, and darwin terminals.
// For windows, it should work fine but no test.
// For others terminals, such as plan9, it's not supported.
func (s *Cursor) CursorGet(ctx context.Context, pos *CursorPos) (chain *Cursor) {
	chain = s

	var n int
	_, _ = fmt.Fprint(os.Stdout, "\033[6n")
	line, ok, err := chk.ReadTill(0, 'R') // read from stdin
	if err != nil {
		slog.Error("getCursorPosTo() readtill failed", "err", err, "n", n, "ok", ok)
		return
	}
	// println("(got):", line[1:], ok, err)
	if line[0] == '\x1b' && line[1] == '[' {
		n, err = fmt.Sscanf(line[2:], "%d;%d", &pos.Row, &pos.Col)
		if err != nil {
			slog.Error("getCursorPosTo() sscanf failed", "err", err, "n", n)
		}
	}
	return
}

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

func writecsi(out Writer, what rune) {
	var bb = []byte(aecHideCursor)
	_, _ = out.Write(bb[0:2])
	_, _ = out.Write([]byte{byte(what)})
}

func writecsiseq(out Writer, what rune, n int) {
	var sb bytes.Buffer
	var bb = []byte(aecHideCursor)
	_, _ = sb.Write(bb[0:2])
	var ss = strconv.Itoa(n)
	_, _ = sb.WriteString(ss)
	_ = sb.WriteByte(byte(what))
	_, _ = out.Write(sb.Bytes())
}

// showCursor shows the cursor.
func showCursor(w Writer) {
	safeWrite(w, []byte(aecShowCursor))
}

// hideCursor hides the cursor.
func hideCursor(w Writer) {
	safeWrite(w, []byte(aecHideCursor))
}

func safeWrite(w Writer, b []byte) (n int, e error) {
	return w.Write(b)
}

// var escape = []byte{'\x1b'}

const aecHideCursor = "\x1b[?25l"
const aecShowCursor = "\x1b[?25h"
