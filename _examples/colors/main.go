package main

import (
	"fmt"

	"github.com/hedzr/is/term/color"
)

func main() {
	run1()
}

func run1() {
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
	fmt.Println(result)

	// another colorful builfer
	c = color.New()
	fmt.Println(c.Color16(color.FgRed).
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

	fmt.Print("0") // now, col = 1
	color.Up(2)
	fmt.Print("ABC") // embedded "ABC" into "[]"
	color.Right(2)   // to be overwrite "hello"
	fmt.Print("HELLO")

	color.RestorePos()
	fmt.Print("z") // write "z" to beginning of "[163]" line

	fmt.Println()
	fmt.Println()
	fmt.Println()

	// color.Down(4)
	// color.Left(1)
	fmt.Println("END")
}
