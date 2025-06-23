package color

import (
	"os"
	"strings"
)

// RowsBlock displays content which can be updated on the fly.
// You can use this to create live output, such as progressbar, etc.
type RowsBlock struct {
	height     int
	writer     Writer
	cursor     *Cursor
	cursorPosY int
}

// NewRowsBlock returns a new RowsBlock.
//
// A RowsBlock displays content which can be updated on the fly.
// You can use this to create live output, such as progressbar, etc.
func NewRowsBlock() RowsBlock {
	return RowsBlock{
		height:     0,
		writer:     os.Stdout,
		cursor:     New(),
		cursorPosY: 0,
	}
}

// WithWriter sets the custom writer.
func (s *RowsBlock) WithWriter(writer Writer) *RowsBlock {
	s.writer = writer
	s.cursor = s.cursor.WithWriter(writer)
	return s
}

// Cursor returns the *cS (Cursor) object so that
// you could render the colorful text with it.
func (s *RowsBlock) Cursor() *Cursor { return s.cursor }

// Clear clears the content of the RowsBlock.
func (s *RowsBlock) Clear() {
	// Initialize writer if not done yet
	if s.writer == nil {
		s.writer = os.Stdout
	}

	if s.height > 0 {
		s.Bottom()
		s.ClearLinesUp(s.height)
		s.Home()
	} else {
		s.Home()
		s.cursor.EraseLineNow()
	}
}

// Update overwrites the content of the RowsBlock
// and adjusts its height based on content.
func (s *RowsBlock) Update(content string) {
	s.Clear()
	s.writeArea(content)
	s.cursorPosY = 0
	s.height = strings.Count(content, "\n")
}

// ShowCursor make the console cursor visible
func (s *RowsBlock) ShowCursor() {
	showCursor(s.writer)
}

// HideCursor make the console cursor invisible
func (s *RowsBlock) HideCursor() {
	hideCursor(s.writer)
}

// Up moves the cursor of the RowsBlock up one line.
func (s *RowsBlock) Up(n int) {
	if n > 0 {
		if s.cursorPosY+n > s.height {
			n = s.height - s.cursorPosY
		}

		s.cursor.UpNow(n)
		s.cursorPosY += n
	}
}

// Down moves the cursor of the RowsBlock down one line.
func (s *RowsBlock) Down(n int) {
	if n > 0 {
		if s.cursorPosY-n < 0 {
			n = s.height - s.cursorPosY
		}

		s.cursor.DownNow(n)
		s.cursorPosY -= n
	}
}

// Bottom moves the cursor to the bottom of the RowsBlock.
// This is done by calculating how many lines were
// moved by Up and Down.
func (s *RowsBlock) Bottom() {
	if s.cursorPosY > 0 {
		s.Down(s.cursorPosY)
		s.cursorPosY = 0
	}
}

// Top moves the cursor to the top of the RowsBlock.
// This is done by calculating how many lines were
// moved by Up and Down.
func (s *RowsBlock) Top() {
	if s.cursorPosY < s.height {
		s.Up(s.height - s.cursorPosY)
		s.cursorPosY = s.height
	}
}

// Home moves the cursor to the start of the current line.
func (s *RowsBlock) Home() {
	s.cursor.HorizontalAbsoluteNow(0)
}

// HomeAndLineDown moves the cursor down by n lines,
// then moves to cursor to the start of the line.
func (s *RowsBlock) HomeAndLineDown(n int) {
	s.Down(n)
	s.Home()
}

// HomeAndLineUp moves the cursor up by n lines, then
// moves to cursor to the start of the line.
func (s *RowsBlock) HomeAndLineUp(n int) {
	s.Up(n)
	s.Home()
}

// UpAndClear moves the cursor up by n lines, then
// clears the line.
func (s *RowsBlock) UpAndClear(n int) {
	s.Up(n)
	s.cursor.EraseLineNow()
}

// DownAndClear moves the cursor down by n lines, then
// clears the line.
func (s *RowsBlock) DownAndClear(n int) {
	s.Down(n)
	s.cursor.EraseLineNow()
}

// Move moves the cursor relative by x and y.
func (s *RowsBlock) Move(x, y int) {
	if x > 0 {
		s.cursor.RightNow(x)
	} else if x < 0 {
		s.cursor.LeftNow(-x)
	}

	if y > 0 {
		s.Up(y)
	} else if y < 0 {
		s.Down(-y)
	}
}

// ClearLinesUp clears n lines upwards from the current
// position and moves the cursor.
func (s *RowsBlock) ClearLinesUp(n int) {
	s.Home()
	s.cursor.EraseLineNow()

	for range n {
		s.UpAndClear(1)
	}
}

// ClearLinesDown clears n lines downwards from the
// current position and moves the cursor.
func (s *RowsBlock) ClearLinesDown(n int) {
	s.Home()
	s.cursor.EraseLineNow()

	for i := 0; i < n; i++ {
		s.DownAndClear(1)
	}
}
