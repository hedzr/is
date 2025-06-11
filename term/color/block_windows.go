//go:build windows
// +build windows

package color

import (
	"fmt"
)

// writeArea is a helper for platform dependant output.
// For Windows newlines '\n' in the content are replaced by '\r\n'
func (s *RowsBlock) writeArea(content string) {
	last := ' '
	for _, r := range content {
		if r == '\n' && last != '\r' {
			fmt.Fprint(s.writer, "\r\n")
			continue
		}
		fmt.Fprint(s.writer, string(r))
	}
}
