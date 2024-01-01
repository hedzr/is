//go:build plan9
// +build plan9

// Copyright © 2020 Hedzr Yeh.

package term

// // ReadPassword reads the password from stdin with safe protection
// func ReadPassword() (text string, err error) {
// 	return stringtool.RandomStringPure(9), nil
// }

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) {
	cols, rows = 32768, 43
	return
}

func isDoubleClickRun() bool { return false }
