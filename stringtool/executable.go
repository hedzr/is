package stringtool

import (
	"os"
	"path/filepath"
)

// GetExecutablePath returns the executable file path
func GetExecutablePath() string {
	p, _ := os.Executable()
	p, _ = filepath.Abs(p)
	return p
}

// GetExecutableDir returns the executable file directory
func GetExecutableDir() string {
	p := GetExecutablePath()
	d := filepath.Dir(p)
	return d
}

// GetCurrentDir returns the current workingFlag directory
func GetCurrentDir() string {
	d, _ := os.Getwd()
	return d
}
