package color

import "io"

func GetDummyTranslator() Translator { return dummy } // return Translator for displaying plain text without color

var _ Translator = (*dummyS)(nil) // ensure cpTranslator implements Translator

var dummy dummyS

type dummyS struct{}

// stripHTMLTags implements Translator.
func (c dummyS) stripHTMLTags(s string) string {
	panic("unimplemented")
}

func (dummyS) Translate(s string, initialFg Color) string        { return s } //nolint:revive
func (dummyS) TranslateTo(s string, initialState Color) string   { return s }
func (dummyS) ColoredFast(out io.Writer, clr Color, text string) { _, _ = out.Write([]byte(text)) } //nolint:revive
func (dummyS) DimFast(out io.Writer, text string)                { _, _ = out.Write([]byte(text)) }
func (dummyS) HighlightFast(out io.Writer, text string)          { _, _ = out.Write([]byte(text)) }
func (dummyS) WriteColor(out io.Writer, clr Color)               {} //nolint:revive
func (dummyS) WriteBgColor(out io.Writer, clr Color)             {} //nolint:revive
func (dummyS) Reset(out io.Writer)                               {} //nolint:revive
func (dummyS) ToColorInt(s string) Color                         { return cptNC.toColorInt(s) }
func (dummyS) ToColorString(clr Color) string                    { return cptNC.toColorString(clr) }

func (c dummyS) StripLeftTabsAndColorize(s string) string { return cptNC.StripLeftTabsAndColorize(s) }
func (c dummyS) StripLeftTabs(s string) string            { return StripLeftTabs(s) }
func (c dummyS) StripLeftTabsOnly(s string) string        { return StripLeftTabsOnly(s) }
func (c dummyS) StripHTMLTags(s string) string            { return StripHTMLTags(s) }

func (c dummyS) Bold(out io.Writer, cb func(out io.Writer))      { c.bg(out, BgBoldOrBright, cb) }
func (c dummyS) Italic(out io.Writer, cb func(out io.Writer))    { c.bg(out, BgItalic, cb) }
func (c dummyS) Underline(out io.Writer, cb func(out io.Writer)) { c.bg(out, BgUnderline, cb) }
func (c dummyS) Inverse(out io.Writer, cb func(out io.Writer))   { c.bg(out, BgInverse, cb) }
func (c dummyS) Dim(out io.Writer, cb func(out io.Writer))       { c.bg(out, BgDim, cb) }
func (c dummyS) Blink(out io.Writer, cb func(out io.Writer))     { c.bg(out, BgBlink, cb) }
func (c dummyS) Bg(out io.Writer, bgColor Color, cb func(out io.Writer)) {
	c.bg(out, bgColor, cb)
}
func (c dummyS) bg(out io.Writer, bg Color, cb func(out io.Writer)) {
	cb(out)
	_ = bg
}
