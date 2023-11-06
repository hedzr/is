package stringtool

import (
	"math/rand"
	"time"
)

// RandomStringPure generate a random string with length specified.
func RandomStringPure(length int) (result string) {
	source := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabets[source.Int63()%int64(len(alphabets))]
	}
	return string(b)
}

const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
