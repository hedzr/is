package is

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hedzr/is/basics"
)

type closerS struct{}

// Closers returns the container includes all registered closable objects.
//
// The simplest ways is using package level Close function:
//
//	func main() {
//	    defer is.Closers().Close()
//
//	    // others statements ...
//
//	    is.Closers().RegisterCloseFns(func(){ sth.Close() })
//
//	    // more statements...
//	}
func Closers() closerS { return closerS{} }

// RegisterPeripheral adds a peripheral into our global closers set.
// a basics.Peripheral object is a closable instance.
func (s closerS) RegisterPeripheral(servers ...basics.Peripheral) {
	basics.RegisterPeripheral(servers...)
}

// RegisterCloseFns adds a simple closure into our global closers set
func (s closerS) RegisterCloseFns(fns ...func()) { basics.RegisterCloseFns(fns...) }

// RegisterClosable adds a peripheral/closable into our global closers set.
// a basics.Peripheral object is a closable instance.
func (s closerS) RegisterClosable(servers ...basics.Closable) { basics.RegisterClosable(servers...) }

// RegisterClosers adds a simple closure into our global closers set
func (s closerS) RegisterClosers(cc ...basics.Closer) { basics.RegisterClosers(cc...) }

// Close will cleanup all registered closers.
// You must make a call to Close before your app shutting down. For example:
//
//	func main() {
//	    defer is.Closers().Close()
//	    // ...
//	}
func (s closerS) Close() {
	basics.Close()
}

// Closers returns the closers set as a basics.Peripheral
func (s closerS) Closers() basics.Peripheral { return basics.Closers() }

// ClosersClosers returns the closers set as a basics.Peripheral array
func (s closerS) ClosersClosers() []basics.Peripheral { return basics.ClosersClosers() }

type signalS struct{}

// Signals returns a signals-helper so that you can catch them, and raise them later.
//
// Typically, its usage is `catcher := is.Signals().Catch(); ...`.
//
// By default, catcher will listen on standard signals set: os.Interrupt,
// os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT.
//
// But you can change them with:
//
//	is.Signals().Catch(os.Kill, os.Interrupt)
//
// or:
//
//	is.Signals().Catch().WithSignals(os.Interrupt, os.Kill)
//
// You should put your long-term codes inside `cb` of WaitFor(cb), and
// defer call to `closer()` in. The `closer()` is a param of `cb`.
//
// For example:
//
//	package main
//
//	import (
//	  "context"
//	  "fmt"
//	  "os"
//	  "sync"
//
//	  "github.com/hedzr/env"
//	  "github.com/hedzr/go-socketlib/net"
//
//	  logz "logslog"
//	)
//
//	func main() {
//	  logz.SetLevel(logz.DebugLevel)
//
//	  server := net.NewServer(":7099")
//	  defer server.Close()
//
//	  ctx, cancel := context.WithCancel(context.Background())
//	  defer cancel()
//
//	  catcher := is.Signals().Catch()
//	  catcher.WithVerboseFn(func(msg string, args ...any) {
//	    logz.Debug(fmt.Sprintf("[verbose] %s", fmt.Sprintf(msg, args...)))
//	  }).
//	  WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
//	    println()
//	    logz.Debug("signal caught", "sig", sig)
//	    if err := server.Shutdown(); err != nil {
//	    	logz.Error("server shutdown error", "err", err)
//	    }
//	    cancel()
//	  }).
//	  WaitFor(func(closer func()) {
//	    logz.Debug("entering looper's loop...")
//
//	    server.WithOnShutdown(func(err error) { closer() })
//	    err := server.StartAndServe(ctx)
//	    if err != nil {
//	      logz.Fatal("server serve failed", "err", err)
//	    }
//	  })
//	}
func Signals() signalS { return signalS{} }

// func (s signalS) SetupCloseHandlerAndWait(wg *sync.WaitGroup, closers ...basics.Peripheral) {
// 	basics.SetupCloseHandlerAndWait(wg, closers...)
// }
//
// func (s signalS) SetupCloseHandlerAndEnterLoop(closers ...basics.Peripheral) {
// 	basics.SetupCloseHandlerAndEnterLoop(closers...)
// }
//
// func (s signalS) SetupCloseHandler(closers ...basics.Peripheral) {
// 	basics.SetupCloseHandler(closers...)
// }

func (s signalS) Catch(sig ...os.Signal) basics.Catcher {
	return basics.Catch(sig...)
}

// RaiseSignal should throw a POSIX signal to current process.
//
// It can work or not, see also the discusses at:
//
//	https://github.com/golang/go/issues/19326
//
// In general cases, it works. But in some special scenes it notifies a quiet thread somewhere.
func (s signalS) RaiseSignal(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func (s signalS) Wait() (*os.ProcessState, error) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return nil, err
	}
	return p.Wait()
}

func (s signalS) Kill() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Kill()
}

func (s signalS) CurrentProcess() *os.Process {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return nil
	}
	return p
}

// PressEnterToContinue lets program pause and wait for user's ENTER key press in console/terminal
func PressEnterToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press 'Enter' to continue...")
	}
	b, _ := bufio.NewReader(in).ReadBytes('\n')
	return strings.TrimRight(string(b), "\n")
}

// PressAnyKeyToContinue lets program pause and wait for user's ANY key press in console/terminal
func PressAnyKeyToContinue(in io.Reader, msg ...string) (input string) {
	if len(msg) > 0 && len(msg[0]) > 0 {
		fmt.Print(msg[0])
	} else {
		fmt.Print("Press any key to continue...")
	}
	_, _ = fmt.Fscanf(in, "%s", &input)
	return
}
