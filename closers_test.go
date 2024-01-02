package is

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hedzr/is/basics"
)

func TestClosers(t *testing.T) {
	defer Closers().Close()

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

	Signals().Catch().
		WithPrompt().
		Wait(func(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup) {
			basics.VerboseFn("[cb] raising interrupt after a second...")
			time.Sleep(2500 * time.Millisecond)
			stopChan <- os.Interrupt
			basics.VerboseFn("[cb] raised.")
		})
}
