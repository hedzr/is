//go:build !windows
// +build !windows

package color

import (
	"fmt"
)

// Update overwrites the content of the Area and adjusts its height based on content.
func (s *RowsBlock) writeArea(content string) {
	fmt.Fprint(s.writer, content)
}
