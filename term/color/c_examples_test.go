package color_test

import (
	"fmt"

	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term/color"
)

func ExampleNew() {
	// start a color text builder
	var c = color.New()

	// specially for running on remote ci server
	if states.Env().IsNoColorMode() {
		states.Env().SetNoColorMode(true)
	}

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

	// For most of ttys, the output looks like:
	//
	// [31mhello, world.[0m
	// [sx
	// [32mhello, world.
	// [38;5;160m[160] hello, world.
	// [38;5;161m[161] hello, world.
	// [38;5;162m[162] hello, world.
	// [38;5;163m[163] hello, world.
	// [38;5;164m[164] hello, world.
	// [38;5;165m[165] hello, world.
	// [0m[3A ERASED [38;2;211;211;33m[16m] hello, world.
	// [uz
	// [8BDONE
}

func ExampleCursor_Color16() {
	// another colorful builfer
	var c = color.New()
	fmt.Println(c.Color16(color.FgRed).
		Printf("hello, %s.", "world").Println().Build())
	// Output: [31mhello, world.[0m
}

func ExampleCursor_Color() {
	// another colorful builfer
	var c = color.New()
	fmt.Println(c.Color(color.FgRed, "hello, %s.", "world").Build())
	// Output: [31mhello, world.[0m
}

func ExampleCursor_Bg() {
	// another colorful builfer
	var c = color.New()
	fmt.Println(c.Bg(color.BgRed, "hello, %s.", "world").Build())
	// Output: [41mhello, world.[0m
}

func ExampleCursor_Effect() {
	// another colorful builfer
	var c = color.New()
	fmt.Println(c.Effect(color.BgDim, "hello, %s.", "world").Build())
	// Output: [2mhello, world.[0m
}

func ExampleCursor_Color256() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c.
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Build())
	// Output:
	// [38;5;163m[163] hello, world.
	// [38;5;164m[164] hello, world.
	// [38;5;165m[165] hello, world.
	// [0m
}

func ExampleCursor_RGB() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c.
		RGB(211, 211, 33).Printf("[16m] hello, %s.\n", "world").
		BgRGB(211, 211, 33).Printf("[16m] hello, %s.\n", "world").
		Build())
	// Output:
	// [38;2;211;211;33m[16m] hello, world.
	// [48;2;211;211;33m[16m] hello, world.
	// [0m
}

func ExampleCursor_EDim() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c. // Color16(color.FgRed).
			EDim("[DIM] hello, %s.\n", "world").String())
	// Output:
	// [2m[DIM] hello, world.
	// [0m
}

func ExampleCursor_Black() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c. // Color16(color.FgRed).
			Black("[BLACK] hello, %s.\n", "world").String())
	// Output:
	// [30m[BLACK] hello, world.
	// [0m
}

func ExampleCursor_BgBlack() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c. // Color16(color.FgRed).
			BgBlack("[BGBLACK] hello, %s.\n", "world").String())
	// Output:
	// [40m[BGBLACK] hello, world.
	// [0m
}

func ExampleCursor_Translate() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c. // Color16(color.FgRed).
			Translate(`<code>code</code> | <kbd>CTRL</kbd>
		<b>bold / strong / em</b>
		<i>italic / cite</i>
		<u>underline</u>
		<mark>inverse mark</mark>
		<del>strike / del </del>
		<font color="green">green text</font>
		`).String())
	// Output:
	// [51;1mcode[0m[39m | [51;1mCTRL[0m[39m
	//		[1mbold / strong / em[0m[39m
	//		[3mitalic / cite[0m[39m
	//		[4munderline[0m[39m
	//		[7minverse mark[0m[39m
	//		[9mstrike / del [0m[39m
	//		[32mgreen text[0m[39m
}

func ExampleCursor_StripLeftTabsColorful() {
	// another colorful builfer
	var c = color.New()
	fmt.Print(c. // Color16(color.FgRed).
			StripLeftTabsColorful(`
		<code>code</code> | <kbd>CTRL</kbd>
		<b>bold / strong / em</b>
		<i>italic / cite</i>
		<u>underline</u>
		<mark>inverse mark</mark>
		<del>strike / del </del>
		<font color="green">green text</font>
		`).String())
	// Output:
	// [51;1mcode[0m[0m | [51;1mCTRL[0m[0m
	// [1mbold / strong / em[0m[0m
	// [3mitalic / cite[0m[0m
	// [4munderline[0m[0m
	// [7minverse mark[0m[0m
	// [9mstrike / del [0m[0m
	// [32mgreen text[0m[0m
}
