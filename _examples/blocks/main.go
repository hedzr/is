package main

import "github.com/hedzr/is/term/color"

func main() {
	runBlock()
}

func runBlock() {
	var blk = color.NewRowsBlock()
	defer blk.Bottom()
	// blk.WithWriter(os.Stdout)
	blk.Update(`
	var result = c.Println().
		Color16(color.FgRed).
		Printf("hello, %s.", "world").Println().
		Color16(color.FgGreen).Printf("hello, %s.\n", "world").
		Color256(160).Printf("[160] hello, %s.\n", "world").
		Color256(161).Printf("[161] hello, %s.\n", "world").
		Color256(162).Printf("[162] hello, %s.\n", "world").
		Color256(163).Printf("[163] hello, %s.\n", "world").
		Color256(164).Printf("[164] hello, %s.\n", "world").
		Color256(165).Printf("[165] hello, %s.\n", "world").
		Build()
`)
	blk.Up(5)
	print("12345")
	blk.Up(2)
	blk.Home()
	print("xyz")
}
