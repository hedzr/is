package color

import (
	"os"
)

// Out is the default output writer for the Writer
var Out Writer = os.Stdout

func Hide() {
	hideCursor(Out)
}

func Show() {
	showCursor(Out)
}

// // Up moves cursor up by n
// func Up(n int) {
// 	syscalls.Up(n)
// }
//
// // Left moves cursor left by n
// func Left(n int) {
// 	syscalls.Left(n)
// }

// Up moves cursor up by n
func Up(n int) { cursorUp(Out, n) }

// Left moves cursor left by n
func Left(n int) { cursorLeft(Out, n) }

func Right(n int) { cursorRight(Out, n) }
func Down(n int)  { cursorLeft(Out, n) }

func ScrollUp(n int)   { cursorScrollUp(Out, n) }
func ScrollDown(n int) { cursorScrollDown(Out, n) }

func SavePos()    { cursorSavePos(Out) }
func RestorePos() { cursorRestorePos(Out) }

func Write(b []byte) (n int, e error) {
	return safeWrite(Out, b)
}

func eraseLine(w Writer, method EraseTo) {
	writecsiseq(w, 'K', int(method))
}

func eraseScreen(w Writer, method EraseTo) {
	writecsiseq(w, 'J', int(method))
}
