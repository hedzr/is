package is

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

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

//

//

//

// SignalsEnh returns a rich-customized struct for operations.
func SignalsEnh() *SignalsX { return &SignalsX{} }

type SignalsX struct {
	dur          time.Duration
	msg          string
	globalCloser func()
	signals      []os.Signal
	cancelled    int32
}

type CatcherOpt func(s *SignalsX)

func WithCatcherMsg(msg string) CatcherOpt {
	return func(s *SignalsX) {
		s.msg = msg
	}
}

func WithCatcherCloser(globalCloser func()) CatcherOpt {
	return func(s *SignalsX) {
		s.globalCloser = globalCloser
	}
}

func WithCatcherDuration(dur time.Duration) CatcherOpt {
	return func(s *SignalsX) {
		s.dur = dur
	}
}

func WithCatcherSignals(sigs ...os.Signal) CatcherOpt {
	return func(s *SignalsX) {
		s.signals = sigs
	}
}

func (s SignalsX) Catch(sig ...os.Signal) basics.Catcher {
	return basics.Catch(sig...)
}

// RaiseSignal should throw a POSIX signal to current process.
//
// It can work or not, see also the discusses at:
//
//	https://github.com/golang/go/issues/19326
//
// In general cases, it works. But in some special scenes it notifies a quiet thread somewhere.
func (s SignalsX) RaiseSignal(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func (s SignalsX) Wait() (*os.ProcessState, error) {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return nil, err
	}
	return p.Wait()
}

func (s SignalsX) Kill() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Kill()
}

func (s SignalsX) CurrentProcess() *os.Process {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return nil
	}
	return p
}

func (s *SignalsX) WaitForContext(ctx context.Context, cancelFunc context.CancelFunc, opts ...CatcherOpt) {
	catcher := s.Catch()

	for _, opt := range opts {
		opt(s)
	}
	if s.msg != "" {
		catcher.WithPrompt(s.msg)
	}
	if len(s.signals) > 0 {
		catcher.WithSignals(s.signals...)
	}

	catcher.
		// WithOnLoopFunc(dbStarter, cacheStarter, mqStarter).
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			if sig != syscall.SIGUSR1 { // not really a signal caught, this means catcher-manager is terminating an onSignalCaught handler.
				slog.Info("signal caught (main)", "sig", sig)
				if s.globalCloser != nil {
					s.globalCloser() // cancel() // cancel user's loop, see <-ctx.Done() in Wait(...)
				}
				if cancelFunc != nil {
					cancelFunc()
				}
			}
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			slog.Debug("entering looper's loop...")
			// go func() {
			defer closer()
			// to terminate this app after a while automatically:
			if s.dur > 0 {
				// time.Sleep(s.dur)
				ticker := time.NewTicker(s.dur)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						return
					case <-ctx.Done():
						return
					}
				}
			}
			// }()
		})
}

func (s *SignalsX) WaitFor(ctx context.Context, opts ...CatcherOpt) {
	s.WaitForContext(ctx, nil, opts...)
}

func (s *SignalsX) WaitForSeconds(ctx context.Context, cancelFunc context.CancelFunc, duration time.Duration, opts ...CatcherOpt) {
	s.dur = duration
	s.WaitForContext(ctx, cancelFunc, opts...)
}

// WaitForSeconds prompts a msg and waits for seconds, or user's pressing CTRL-C to quit.
//
//	package main
//
//	import (
//		"context"
//		"time"
//
//		"github.com/hedzr/is"
//		"github.com/hedzr/is/timing"
//	)
//
//	func main() {
//		ctx, cancel := context.WithCancel(context.Background())
//		defer cancel()
//
//		p := timing.New()
//		defer p.CalcNow()
//
//		is.SignalsEnh().WaitForSeconds(ctx, cancel, 6*time.Second,
//			// is.WithCatcherCloser(cancel),
//			is.WithCatcherMsg("Press CTRL-C to quit, or waiting for 6s..."),
//		)
//	}
func WaitForSeconds(ctx context.Context, cancelFunc context.CancelFunc, duration time.Duration, opts ...CatcherOpt) {
	SignalsEnh().WaitForSeconds(ctx, cancelFunc, duration, opts...)
}

//

//

//

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
