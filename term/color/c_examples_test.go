package color_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
		SavePos(). // try using SavePosNow()
		Println("x").
		Color16(color.FgGreen).Printf("hello, %s.\n", "world").
		Color256(160).Printf("[160] hello, %s.\n", "world").
		Color256(161).Printf("[161] hello, %s.\n", "world").
		Color256(162).Printf("[162] hello, %s.\n", "world").
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Up(3).Echo(" ERASED "). // try using UpNow()
		RGB(211, 211, 33).Printf("[16m] hello, %s.", "world").
		Println().
		RestorePos(). // try using RestorePosNow()
		Println("z").
		Down(8). // try using DownNow()
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

func ExampleNewColor16() {
	var clr = color.NewColor16(color.FgCyan)
	var str = clr.Color()
	fmt.Println(str, "hello")
	// Output:
	// [36m hello
}

func ExampleNewColor256() {
	var clr = color.NewColor256(byte(137), false)
	fmt.Println(clr.Color(), "hello")
	clr = color.NewColor256(byte(137), true)
	fmt.Println(clr.Color(), "hello")
	// Output:
	// [38;5;137m hello
	// [48;5;137m hello
}

func ExampleNewColor16m() {
	var clr = color.NewColor16m(173, 137, 73, false)
	fmt.Println(clr.Color(), "hello")
	clr = color.NewColor16m(173, 137, 73, true)
	fmt.Println(clr.Color(), "hello")
	// Output:
	// [38;2;173;137;73m hello
	// [48;2;173;137;73m hello
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

func ExampleCSICodes() {
	c := color.CSICursorUp.Code1(7)
	fmt.Printf("%s", c) // move cursor up 7 lines
	d := color.CSICursorDown.Code1(7)
	fmt.Printf("%s", d) // move cursor down 7 lines

	// move cursor to new position

	fmt.Printf("%s", color.New().SavePosNow().Build())

	var cx color.CSICodes
	cx.AddCode2(color.CSICursorPosition, 2, 3) // move cursor to row 2 col 3
	fmt.Printf("%s", cx)

	fmt.Printf("%s", color.New().RestorePosNow().Build())

	// Output s:
	// [7A[7B[s[2,3H[u
}

func ExampleNewStyle() {
	c := color.NewStyle()
	c.Add(
		color.NewColor16(color.FgYellow),    // fg
		color.NewColor16m(77, 88, 99, true), // bg
		// color.Reset,
	)
	fmt.Printf("%sHello, World!%s\n", c, color.Reset)
	// Output:
	// [33m[48;2;77;88;99mHello, World![0m
}

func ExampleNewRowsBlock() {
	rb := color.NewRowsBlock()

	// the following outputs will be displayed in first
	// line of the RowsBlock.
	for ul := range 10 {
		spc := strings.Repeat("+", ul)
		str := fmt.Sprintf("%sHello, World!\n", spc)
		rb.Update(str)
	}

	// don't test this example because the outputs on different tty (or ci servers) could fail.

	// Outputs:
	// [0G[2KHello, World!
	// [0G[2K[1A[2K[0G+Hello, World!
	// [0G[2K[1A[2K[0G++Hello, World!
	// [0G[2K[1A[2K[0G+++Hello, World!
	// [0G[2K[1A[2K[0G++++Hello, World!
	// [0G[2K[1A[2K[0G+++++Hello, World!
	// [0G[2K[1A[2K[0G++++++Hello, World!
	// [0G[2K[1A[2K[0G+++++++Hello, World!
	// [0G[2K[1A[2K[0G++++++++Hello, World!
	// [0G[2K[1A[2K[0G+++++++++Hello, World!
}

func TestExampleNewRowsBlock(t *testing.T) {
	if !testing.Verbose() {
		return
	}

	rb := color.NewRowsBlock()

	// the following outputs will be displayed in first
	// line of the RowsBlock.
	for ul := range 10 {
		spc := strings.Repeat("+", ul)
		str := fmt.Sprintf("%sHello, World!\n", spc)
		rb.Update(str)
	}
}

func ExampleNewSGR() {
	for i, sgrs := range []struct {
		pre, post color.CSIsgr
		desc      string
	}{
		{color.SGRbold, color.SGRresetBoldAndDim, "bold"},
		{color.SGRdim, color.SGRresetBoldAndDim, "dim"},
		{color.SGRitalic, color.SGRresetItalic, "italic"},
		{color.SGRunderline, color.SGRresetUnderline, "underline"},
		{color.SGRslowblink, color.SGRresetSlowBlink, "blink"},
		{color.SGRrapidblink, color.SGRresetRapidBlink, "fast blink"},
		{color.SGRinverse, color.SGRresetInverse, "inverse"},
		{color.SGRhide, color.SGRresetHide, "hide"},
		{color.SGRstrike, color.SGRresetStrike, "strike"},
		{color.SGRframed, color.SGRneitherFramedNorEncircled, "framed"},
		{color.SGRencircled, color.SGRneitherFramedNorEncircled, "encircled"},
		{color.SGRoverlined, color.SGRnotoverlined, "overlined"},
		{color.SGRideogramUnderline, color.SGRresetIdeogram, "ideogram underline"},
		{color.SGRideogramDoubleUnderline, color.SGRresetIdeogram, "ideogram double underline"},
		{color.SGRideogramOverline, color.SGRresetIdeogram, "ideogram overline"},
		{color.SGRideogramDoubleOverline, color.SGRresetIdeogram, "ideogram double overline"},
		{color.SGRideogramStressMarking, color.SGRresetIdeogram, "ideogram stress marking"},
		{color.SGRsuperscript, color.SGRresetSuperscriptAndSubscript, "superscript"},
		{color.SGRsubscript, color.SGRresetSuperscriptAndSubscript, "subscript"},
		// {color.SGRdim, color.SGRresetDim},
	} {
		str := fmt.Sprintf(`%5d. %s%s%s %s`,
			i, sgrs.pre,
			"Hello, World!",
			sgrs.post,
			sgrs.desc,
		)
		fmt.Println(str)
	}

	fmt.Println(color.SGRreset)

	// Output:
	//     0. [1mHello, World![22m bold
	//     1. [2mHello, World![22m dim
	//     2. [3mHello, World![23m italic
	//     3. [4mHello, World![24m underline
	//     4. [5mHello, World![25m blink
	//     5. [6mHello, World![26m fast blink
	//     6. [7mHello, World![27m inverse
	//     7. [8mHello, World![28m hide
	//     8. [9mHello, World![29m strike
	//     9. [51mHello, World![54m framed
	//    10. [52mHello, World![54m encircled
	//    11. [53mHello, World![55m overlined
	//    12. [60mHello, World![65m ideogram underline
	//    13. [61mHello, World![65m ideogram double underline
	//    14. [62mHello, World![65m ideogram overline
	//    15. [63mHello, World![65m ideogram double overline
	//    16. [64mHello, World![65m ideogram stress marking
	//    17. [73mHello, World![75m superscript
	//    18. [74mHello, World![75m subscript
	// [0m
}

func TestExampleNewSGR(t *testing.T) {
	if !testing.Verbose() {
		return
	}
	ExampleNewSGR()

	// SGRsetFg
	fmt.Printf("\x1b[%d;5;9m[ 9 TEST string HERE]%s\n",
		color.SGRsetFg,
		color.SGRdefaultFg,
	)

	fmt.Printf("\x1b[%d;5;21m[21 TEST string HERE]%s\n",
		color.SGRsetFg,
		color.SGRdefaultFg,
	)

	fmt.Println(color.SGRreset)

	// [38;5;9m[ 9 TEST string HERE][39m
	// [38;5;21m[21 TEST string HERE][39m
	// [0m
}

func TestColor256table(t *testing.T) {
	color.Color256table(os.Stdout)
}
