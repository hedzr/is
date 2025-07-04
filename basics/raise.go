package basics

import (
	"os"
)

func Raise(sig os.Signal) error {
	return raiseOsSig(sig)
}
