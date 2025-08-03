package is

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/hedzr/is/basics"
)

func TestClosers(t *testing.T) {
	defer Closers().Close()

	t.Log(DebuggerAttached())
}

type redisHub struct{}

func (s *redisHub) Close() {
	// close the connections to redis servers
	println("redis connections closed")
}

func TestCloserStruct_Close(t *testing.T) {
	defer func() { println("closed.") }()
	defer Closers().Close()
	defer func() { println("closers.Close() will be invoked at program terminating.") }()

	t.Log("running")

	Closers().RegisterCloseFns(func() {
		// do some shutdown operations here
		println("close functor")
	})
	Closers().RegisterPeripheral(&redisHub{})

	tmpFile, err := os.CreateTemp(os.TempDir(), "1*.log")
	t.Logf("tmpfile: %v | err: %v", tmpFile.Name(), err)
	basics.RegisterClosers(tmpFile)
}

func TestSignalStruct_Raise(t *testing.T) {
	// not a true test here

	basics.VerboseFn = t.Logf

	// done := make(chan struct{})
	go func() {
		basics.VerboseFn("go routine 1 started.")
		time.Sleep(200 * time.Millisecond)
		basics.VerboseFn("go routine 1 stopped.")
	}()

	ctx := context.Background()
	Signals().Catch().
		WithPrompt().
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			go func() {
				basics.VerboseFn("[cb] raising interrupt after a second...")
				time.Sleep(2500 * time.Millisecond)
				basics.VerboseFn("[cb] raised.")
				closer()
			}()
		})
}
