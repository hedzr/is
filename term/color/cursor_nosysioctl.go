//go:build plan9 || appengine || wasm
// +build plan9 appengine wasm

package color

import (
	"os"
)

// Copyright © 2022 Atonal Authors
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
