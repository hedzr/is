package main

import (
	"fmt"

	"github.com/hedzr/is"
	"github.com/hedzr/is/term/color"
)

func main() {
	println(is.InTesting())
	println(is.Env().GetDebugLevel())
	fmt.Printf("%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
	<b>bold / strong / em</b>
	<i>italic / cite</i>
	<u>underline</u>
	<mark>inverse mark</mark>
	<del>strike / del </del>
	<font color="green">green text</font>
	`, color.FgDefault))
}
