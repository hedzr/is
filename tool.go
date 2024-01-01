package is

import (
	"os"

	"github.com/hedzr/is/term"
)

// func RandomStringPure(length int) (result string) {
// 	buf := make([]byte, length)
// 	if _, err := rand.Read(buf); err == nil { //nolint:gosec //like it
// 		result = string(buf)
// 	}
// 	return
// }

// // RandomStringPure generate a random string with length specified.
// func RandomStringPure(length int) (result string) {
// 	source := rand.NewSource(time.Now().UnixNano())
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = Alphabets[source.Int63()%int64(len(Alphabets))]
// 	}
// 	return string(b)
// }

// FileExists detects if a file or a directory is existed.
func FileExists(filepath string) bool { return fileExists(filepath) }

// fileExists returns the existence of an directory or file
func fileExists(filepath string) bool {
	if _, err := os.Stat(os.ExpandEnv(filepath)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ReadFile reads the file named by filename and returns the contents.
func ReadFile(filename string) ([]byte, error) { return readFile(filename) }

// readFile reads the file named by filename and returns the contents.
//
// A successful call returns err == nil, not err == EOF. Because ReadFile
// reads the whole file, it does not treat an EOF from Read as an error
// to be reported.
//
// As of Go 1.16, this function simply calls os.ReadFile.
func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// ToBool converts a value (int, bool, string) to boolean
func ToBool(val any, defaultVal ...bool) (ret bool) {
	return term.ToBool(val, defaultVal...)
}

// StringToBool converts a string to boolean value.
func StringToBool(val string, defaultVal ...bool) (ret bool) {
	return term.StringToBool(val, defaultVal...)
}
