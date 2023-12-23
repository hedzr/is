package basics_test

import (
	"os"
	"testing"

	"github.com/hedzr/is/basics"
)

type redisHub struct{}

func (s *redisHub) Close() {
	// close the connections to redis servers
	println("redis connections closed")
}

func TestClosers(t *testing.T) {
	defer basics.Close()

	basics.RegisterPeripheral(&redisHub{})

	basics.RegisterCloseFns(func() {
		// do some shutdown operations here
		println("close functor")
	})

	basics.RegisterCloseFn(func() {
		// do some shutdown operations here
		println("close single functor")
	})

	tmpFile, err := os.CreateTemp(os.TempDir(), "1*.log")
	t.Logf("tmpfile: %v | err: %v", tmpFile.Name(), err)
	basics.RegisterClosers(tmpFile)

	for _, ii := range basics.ClosersClosers() {
		println(ii)
	}

	// These following calls are both unused since we
	// have had a defer basics.Close().
	// But they are harmless here.

	basics.Closers().Close() //
	basics.Close()           // rerun Close() is safe
}
