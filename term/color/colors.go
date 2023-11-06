// Package color provides a wrapped standard output device like printf but with colored enhancements.
package color

type Color int // ANSI Escaped Sequences here

const (
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	// https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97

	// FgBlack terminal color code
	FgBlack Color = 30
	// FgRed terminal color code
	FgRed Color = 31
	// FgGreen terminal color code
	FgGreen Color = 32
	// FgYellow terminal color code
	FgYellow Color = 33
	// FgBlue terminal color code
	FgBlue Color = 34
	// FgMagenta terminal color code
	FgMagenta Color = 35
	// FgCyan terminal color code
	FgCyan Color = 36
	// FgLightGray terminal color code (White)
	FgLightGray Color = 37

	// FgDarkGray terminal color code (Gray, Light Black).
	//
	// A highlight/bright black color, maybe 50% gray. See FgLightBlack.
	FgDarkGray Color = 90
	// FgLightBlack terminal color code (Gray, Light Black).
	//
	// A highlight/bright black color, maybe 50% gray. See FgDarkGray.
	FgLightBlack Color = 90
	// FgLightRed terminal color code
	FgLightRed Color = 91
	// FgLightGreen terminal color code
	FgLightGreen Color = 92
	// FgLightYellow terminal color code
	FgLightYellow Color = 93
	// FgLightBlue terminal color code
	FgLightBlue Color = 94
	// FgLightMagenta terminal color code
	FgLightMagenta Color = 95
	// FgLightCyan terminal color code
	FgLightCyan Color = 96
	// FgWhite terminal color code (Light White)
	FgWhite Color = 97

	// BgBlack terminal color code
	BgBlack Color = 40
	// BgRed terminal color code
	BgRed Color = 41
	// BgGreen terminal color code
	BgGreen Color = 42
	// BgYellow terminal color code
	BgYellow Color = 43
	// BgBlue terminal color code
	BgBlue Color = 44
	// BgMagenta terminal color code
	BgMagenta Color = 45
	// BgCyan terminal color code
	BgCyan Color = 46
	// BgLightGray terminal color code
	BgLightGray Color = 47
	// BgDarkGray terminal color code
	BgDarkGray Color = 100
	// BgLightRed terminal color code
	BgLightRed Color = 101
	// BgLightGreen terminal color code
	BgLightGreen Color = 102
	// BgLightYellow terminal color code
	BgLightYellow Color = 103
	// BgLightBlue terminal color code
	BgLightBlue Color = 104
	// BgLightMagenta terminal color code
	BgLightMagenta Color = 105
	// BgLightCyan terminal color code
	BgLightCyan Color = 106
	// BgWhite terminal color code
	BgWhite Color = 107

	// BgNormal terminal color code.
	//
	// All attributes become turned off.
	BgNormal Color = 0
	// BgBoldOrBright terminal color code
	//
	// Bold or increased intensity
	BgBoldOrBright Color = 1
	// BgDim terminal color code.
	//
	// Faint, decreased intensity, or dim.
	// May be implemented as a light font weight like bold.
	BgDim Color = 2
	// BgItalic terminal color code.
	//
	// Not widely supported. Sometimes treated as inverse or blink
	BgItalic Color = 3
	// BgUnderline terminal color code.
	//
	// Style extensions exist for Kitty, VTE, mintty, iTerm2 and Konsole.
	BgUnderline Color = 4
	// BgBlink terminal color code.
	//
	// Slow blink, Sets blinking to less than 150 times per minute.
	// But in many tty it's no effect.
	//
	// Sometimes it can be used for switching to 'normal' bg state without
	// reset all fg and bg settings (if using bg code 0)
	BgBlink Color = 5
	// BgRapidBlink terminal color code.
	//
	// MS-DOS ANSI.SYS, 150+ per minute; not widely supported.
	//
	// Sometimes it can be used for switching to 'normal' bg state without
	// reset all fg and bg settings (if using bg code 0)
	BgRapidBlink Color = 6
	// BgInverse terminal color code.
	//
	// Swap foreground and background colors; inconsistent emulation
	BgInverse Color = 7
	// BgHidden terminal color code.
	//
	// not widely supported.
	BgHidden Color = 8
	// BgStrikeout terminal color code.
	//
	// marked as if for deletion.
	BgStrikeout Color = 9

	FgDarkColor = FgLightGray

	FgDefault Color = 39
	BgDefault Color = 49

	// NoColor is not a declared ansi code but we can use it for identifying
	// a variable isn't initializing yet.
	NoColor Color = -1
)
