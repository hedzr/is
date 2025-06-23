package main

import (
	"context"
	"fmt"

	"github.com/hedzr/is/term/color"
)

func main() { run1() }
func run1() {
	// start a color text builder
	var c = color.New()
	var pos color.CursorPos
	ctx := context.Background() // or with cancel

	// paint and get the result (with ansi-color-seq ready)
	var result = c.Println().
		Color16(color.FgRed).Printf("[1st] hello, %s.", "world").
		Println().
		SavePosNow().
		Println("XX").
		Color16(color.FgGreen).Printf("hello, %s.\n", "world").
		Color256(160).Printf("[160] hello, %s.\n", "world").
		Color256(161).Printf("[161] hello, %s.\n", "world").
		Color256(162).Printf("[162] hello, %s.\n", "world").
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		UpNow(4).Echo(" ERASED ").
		RightNow(11).
		CursorGet(ctx, &pos).
		RGB(211, 211, 33).Printf("[16m] hello, %s. pos=%+v", "world", pos).
		Println().
		RestorePosNow().
		Println("ZZ").
		DownNow(8).
		Println("DONE").
		Build()

		// and render the result
	fmt.Println(result)

	// another colorful builfer
	c = color.New()
	fmt.Println(c.Color16(color.FgRed).
		Printf("[2nd] hello, %s.", "world").Println().Build())

	// cursor operations
	c = color.New()
	c.SavePosNow()
	// fmt.Println(c.CursorSavePos().Build())

	fmt.Print(c.
		Printf("[3rd] hello, %s.", "world").
		Println().
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Build())

	fmt.Print("0")         // now, col = 1
	c.UpNow(2)             //
	fmt.Print("ABC")       // embedded "ABC" into "[]"
	c.CursorGet(ctx, &pos) //
	c.RightNow(2)          // to be overwrite "hello"
	fmt.Print("HELLO")     //

	c.RestorePosNow()
	c.DownNow(1)
	fmt.Print("T") // write "T" to beginning of "[163]" line

	c.DownNow(4)

	// color.Down(4)
	// color.Left(1)
	fmt.Printf("\nEND (pos = %+v)\n", pos)
}
