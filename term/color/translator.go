package color

import "io"

// Translator _
type Translator interface {
	Translate(s string, initialFg Color) string

	// StripLeftTabsAndColorize(s string) string

	ColoredFast(out io.Writer, clr Color, text string)
	DimFast(out io.Writer, text string)
	HighlightFast(out io.Writer, text string)

	WriteColor(out io.Writer, clr Color)   // echo ansi color bytes for foreground color
	WriteBgColor(out io.Writer, clr Color) // echo ansi color bytes for background color
	Reset(out io.Writer)                   // echo ansi color bytes for resetting color

	Bold(out io.Writer, cb func(out io.Writer))
	Italic(out io.Writer, cb func(out io.Writer))
	Underline(out io.Writer, cb func(out io.Writer))
	Inverse(out io.Writer, cb func(out io.Writer))
	Dim(out io.Writer, cb func(out io.Writer))
	Blink(out io.Writer, cb func(out io.Writer))
	Bg(out io.Writer, bgColor Color, cb func(out io.Writer))

	// TranslateTo(s string, initialState Color) string
	TranslateTo(s string, initialState Color) string // translate a string with html tags to a colored string

	StripLeftTabsAndColorize(s string) string // strip left tabs and colorize the string
	StripLeftTabs(s string) string            // strip left tabs and colorize the string
	StripLeftTabsOnly(s string) string        // strip left tabs only

	stripHTMLTags(s string) string // aggressively strip HTML tags from a string
	ToColorString(clr Color) string
	ToColorInt(s string) Color
}
