package color_test

import (
	"fmt"
	"testing"

	"github.com/hedzr/is/term/color"
)

func TestColor16(t *testing.T) {
	var c = color.New()

	c.Color(color.FgLightRed, "hello, %s.", "world") // or c.Cyan("hi"), ...
	c.Bg(color.BgCyan, "hello, %s.", "world")        // or c.BgBlack("hi"), ...
	c.Effect(color.BgDim, "hello, %s.", "world")     // or c.EDim("hi"), ...

	c.Println()
	c.Color16(color.FgRed).Printf("hello, %s.", "world").Println() // don't close it till Reset/Println

	// don;t close, but String() will close it automatically
	c.Color16(color.FgGreen).Printf("hello, %s.", "world\n")

	t.Logf("%s", c.String())
}

func TestCSI(t *testing.T) {
	var c = color.New()

	c.Println()
	// up 3 line
	c.CSI('A', 3).Printf("[up] hello, %s.", "world").Println()
	c.CSI('B', 3).Printf("[dn] hello, %s.", "world").Println()

	c.Printf("[BK] hello, %s.", "world").
		CursorBack(6).Echo("there").
		Println()

	fmt.Print(c.String())
	t.Log("OK")
}

func TestColors(t *testing.T) {
	// var c = color.New()
	//
	// c.Println()
	// c.Color16(color.FgRed).
	// 	Printf("hello, %s.", "world").Println().
	// 	CursorSavePos().
	// 	Println("x").
	// 	Color16(color.FgGreen).Printf("hello, %s.\n", "world").
	// 	Color256(160).Printf("[160] hello, %s.\n", "world").
	// 	Color256(161).Printf("[161] hello, %s.\n", "world").
	// 	Color256(162).Printf("[162] hello, %s.\n", "world").
	// 	Color256(163).Printf("[163] hello, %s.\n", "world").
	// 	Color256(164).Printf("[164] hello, %s.\n", "world").
	// 	Color256(165).Printf("[165] hello, %s.\n", "world").
	// 	RGB(211, 211, 33).Printf("[16m] hello, %s.", "world").
	// 	Println().
	// 	CursorRestorePos().
	// 	Println("z").
	// 	CursorDown(8).
	// 	Println("END")
	//
	// t.Logf("%s", c.String())

	// start a color text builder
	var c = color.New()

	// paint and get the result (with ansi-color-seq ready)
	var result = c.Println().
		Color16(color.FgRed).
		Printf("hello, %s.", "world").Println().
		SavePos().
		Println("x").
		Color16(color.FgGreen).Printf("hello, %s.\n", "world").
		Color256(160).Printf("[160] hello, %s.\n", "world").
		Color256(161).Printf("[161] hello, %s.\n", "world").
		Color256(162).Printf("[162] hello, %s.\n", "world").
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Up(3).Echo(" ERASED ").
		RGB(211, 211, 33).Printf("[16m] hello, %s.", "world").
		Println().
		RestorePos().
		Println("z").
		Down(8).
		Println("DONE").
		Build()

		// and render the result
	// fmt.Println(result)
	t.Logf("%s", result)

	// another colorful builfer
	c = color.New()
	t.Logf("%s", c.Color16(color.FgRed).
		Printf("hello, %s.", "world").Println().Build())

	// cursor operations
	c = color.New()
	color.SavePos()
	// fmt.Println(c.CursorSavePos().Build())

	fmt.Print(c.
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Build())

	t.Log("0") // now, col = 1
	color.Up(2)
	t.Log("ABC")   // embedded "ABC" into "[]"
	color.Right(2) // to be overwrite "hello"
	t.Log("HELLO")

	color.RestorePos()
	t.Log("z") // write "z" to beginning of "[163]" line

	fmt.Println()
	fmt.Println()
	fmt.Println()

	// color.Down(4)
	// color.Left(1)
	t.Log("END")
}
