package basics

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// func SetupCloseHandlerAndWait(wg *sync.WaitGroup, closers ...Peripheral) {
// 	setupCloseHandler1(closers...)
// 	wg.Wait()
// }
//
// func SetupCloseHandlerAndEnterLoop(closers ...Peripheral) {
// 	enterLoop(setupCloseHandler1(closers...))
// }
//
// func SetupCloseHandler(closers ...Peripheral) {
// 	setupCloseHandler1(closers...)
// }
//
// func setupCloseHandler1(onFinish ...Peripheral) chan struct{} {
// 	return setupCloseHandlers([]os.Signal{os.Interrupt, syscall.SIGTERM}, onFinish...)
// }
//
// // setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// // program if it receives an interrupt from the OS. We then handle this by calling
// // our clean up procedure and exiting the program.
// func setupCloseHandlers(signals []os.Signal, onFinish ...Peripheral) chan struct{} {
// 	c := make(chan os.Signal, 16)
// 	signal.Notify(c, signals...) //nolint:govet
// 	done := make(chan struct{})
// 	go func() {
// 		<-c
// 		fmt.Println("\r- Ctrl+C pressed in Terminal")
// 		for _, f := range onFinish {
// 			f.Close()
// 		}
// 		closers.Close()
// 		// os.Exit(0)
// 		close(done)
// 	}()
// 	return done
// }
//
// func enterLoop(done chan struct{}) {
// 	for { //nolint:gosimple
// 		select {
// 		case <-done:
// 			return
// 		}
// 	}
// }

// Catch returns a builder to build the programmatic structure for entering a
// infinite loop and waiting for os signals caught or trigger anyone of them
// programmatically.
//
// At the ending of program, all closers (see Peripheral and Close) will be
// closed safely, except panic in their Close codes.
//
// For example,
//
//	basics.VerboseFn = t.Logf
//	is.Signals().Catch().
//	  WithPrompt().
//	  WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
//	    println()
//	    slog.Info("signal caught", "sig", sig)
//	    cancel() // cancel user's loop, see Wait(...)
//	  }).
//	  WaitFor(func(closer func()) {
//	    go func() {
//	      defer closer()
//	      basics.VerboseFn("[cb] raising interrupt after a second...")
//	      time.Sleep(2500 * time.Millisecond)
//	      <-ctx.Done() // waiting for main program stop.
//	      basics.VerboseFn("[cb] raised.")
//	    }()
//	  })
//
// A simple details can be found at:
//
//	https://www.developer.com/languages/os-signals-go/
//
// Your logic that shutdown the main loop gracefully could be:
//
//	type appS struct{}
//	func (s *appS) MainRunner(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup) {
//	     wgShutdown.Done()
//	     stopChan <- os.Interrupt
//	}
//
//	var app appS
//
//	env.Signals().Catch().
//	    Wait(app.MainRunner)
func Catch(signals ...os.Signal) Catcher {
	return &catsig{signals: signals}
}

// Catcher is a builder to build the programmatic structure for entering a
// infinite loop and waiting for os signals caught or trigger anyone of them
// programmatically.
//
// At the ending of program, all closers (see Peripheral and Close) will be
// closed safely, except panic in their Close codes.
type Catcher interface {
	// WithSignals declares a set of signals which can be caught with
	// application processing logics.
	//
	// An empty call means the default set will be applying:
	//
	//	signals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP}
	//
	// These signals are hard-coded.
	WithSignals(signals ...os.Signal) Catcher
	// WithCloseHandlers gives the extra Peripheral's a chance
	// with safely shutting down.
	//
	// The better way is registering them with env.Closers().RegisterPeripheral(p).
	// All registered peripherals will be released/closed in the ending of Wait.
	WithCloseHandlers(onFinish ...Peripheral) Catcher
	// WithPrompt show a message before entering main loop.
	//
	// While many messages given, the final one will be used.
	// Use an empty call like WithPrompt() to request a default
	// prompt line.
	//
	// If you dislike printing anything, do Wait() directly
	// without WithPrompt call.
	WithPrompt(msg ...string) Catcher
	// WithOnSignalCaught setups handlers while any os signals caught by app.
	WithOnSignalCaught(cb ...OnSignalCaught) Catcher
	// WithOnLoopFunc _
	// Deprecated since v0.8.x
	WithOnLoopFunc(cb ...OnLooperFunc) Catcher
	// WithPeripherals _
	WithPeripherals(p ...Peripheral) Catcher
	// WithVerboseFn gives a change to log he catcher's internal state.
	WithVerboseFn(cb func(msg string, args ...any)) Catcher
	// WaitFor with param `cb func(closer func())` is used for your task.
	//
	// You should put your long-term codes inside `cb` of WaitFor(cb), and
	// defer call to `closer()` in. The `closer()` is a param of `cb`.
	//
	// A sample code is,
	//
	//	WaitFor(func(closer func()) {
	//	  defer closer()
	//	  for{
	//	  case <-ticker.C:
	//	    wakeupForTask()
	//	  case <-ctx.Done():
	//	    return
	//	  }
	//	})
	//
	// A http server could be:
	//
	//	WaitFor(func(closer func()) {
	//	   logz.Debug("entering looper's loop...")
	//
	//	   server.WithOnShutdown(func(err error, ss net.Server) { closer() })
	//	   err := server.ListenAndServe(ctx, nil)
	//	   if err != nil {
	//	     logz.Fatal("server serve failed", "err", err)
	//	   }
	//	})
	WaitFor(ctx context.Context, mainLooper OnLooperFunc)
}

type OnSignalCaught func(ctx context.Context, sig os.Signal, wgShutdown *sync.WaitGroup) // callback while an OS signal caught
type OnLooper func(stopChan chan<- os.Signal, wgDone *sync.WaitGroup)                    // callback while get into waiting loop
type OnLooperFunc func(ctx context.Context, closer func())                               // callback while get into waiting loop

type catsig struct {
	signals         []os.Signal
	closers         c
	onCaught        []OnSignalCaught
	looperHandlers  []OnLooper
	looperHandlers1 []OnLooperFunc
	openPeripherals []func(ctx context.Context) (err error)
	msg             string
}

func (s *catsig) Close() {
	verbose("closing catsig.closers...")
	s.closers.Close()
}

func (s *catsig) WithCloseHandlers(onFinish ...Peripheral) Catcher {
	s.closers.RegisterPeripheral(onFinish...)
	return s
}

func (s *catsig) WithClosable(onFinish ...Closable) Catcher {
	s.closers.RegisterClosable(onFinish...)
	return s
}

func (s *catsig) WithSignals(signals ...os.Signal) Catcher {
	s.signals = append(s.signals, signals...)
	return s
}

func (s *catsig) WithPrompt(msg ...string) Catcher {
	var text string
	for _, str := range msg {
		if str != "" {
			text = str
		}
	}
	if text == "" {
		text = "\r- Ctrl+C pressed in Terminal"
	}
	s.msg = text
	return s
}

func (s *catsig) WithOnSignalCaught(cb ...OnSignalCaught) Catcher {
	s.onCaught = cb
	return s
}

func (s *catsig) WithOnLoop(cb ...OnLooper) Catcher {
	s.looperHandlers = append(s.looperHandlers, cb...)
	return s
}

func (s *catsig) WithOnLoopFunc(cb ...OnLooperFunc) Catcher {
	s.looperHandlers1 = append(s.looperHandlers1, cb...)
	return s
}

func (s *catsig) WithPeripherals(peripherals ...Peripheral) Catcher {
	for _, peripheral := range peripherals {
		RegisterPeripheral(peripheral)
		if p, ok := peripheral.(Openable); ok {
			s.openPeripherals = append(s.openPeripherals, p.Open)
		}
	}
	return s
}

func (s *catsig) WithVerboseFn(cb func(msg string, args ...any)) Catcher {
	VerboseFn = cb
	return s
}

func (s *catsig) WaitFor(ctx context.Context, mainLooper OnLooperFunc) {
	defer Close()
	defer s.Close()

	var wgInitialized sync.WaitGroup
	var wgForShutdown sync.WaitGroup
	var closed int32

	done := make(chan struct{}, 8)
	shutDone := func() {
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			verbose("closing done chan finally...")
			close(done)
		}
	}
	defer shutDone()

	cc := make(chan os.Signal, 16) // allow buffering some signals
	signals := s.signals

	var looperHandlers []OnLooperFunc
	looperHandlers = append(looperHandlers, s.looperHandlers1...)
	count := len(looperHandlers) + 1
	wgInitialized.Add(count)
	wgForShutdown.Add(count)
	verbose("...%d looper functions added\n", count)
	for _, cb := range s.looperHandlers {
		looperHandlers = append(looperHandlers, func(ctx context.Context, closer func()) {
			defer closer()
			cb(cc, &wgForShutdown)
		})
	}
	looperHandlers = append(looperHandlers, mainLooper)

	if ctx == nil {
		ctx = context.Background()
	}
	for _, p := range s.openPeripherals {
		p(ctx)
	}

	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGUSR1}
	}
	signal.Notify(cc, signals...) //nolint:govet //whyNoLint for why

	for _, f := range looperHandlers {
		go func(cc chan os.Signal, wgInitialized, wgForShutdown *sync.WaitGroup, f OnLooperFunc) {
			wgInitialized.Done()
			f(ctx, func() { cc <- syscall.SIGUSR1; wgForShutdown.Done() })
		}(cc, &wgInitialized, &wgForShutdown, f)
	}
	wgInitialized.Wait()

	verbose("all looper(s) ran\n")

	go func(wg *sync.WaitGroup) {
		verbose("waiting for os signals...")
		sig := <-cc
		// verbose("caught a signal.")
		for _, cb := range s.onCaught {
			if cb != nil {
				cb(ctx, sig, wg)
			}
		}
		if len(s.onCaught) == 0 {
			println()
			verbose("signal caught", "sig", sig)
		}
		// verbose("wgForShutdown WAIT...\n")
		wg.Wait()
		verbose("wgForShutdown DONE.\n")
		done <- struct{}{}
	}(&wgForShutdown)

	// enter the main loop here till someone raises a signal from looperHandlers
	// by triggering such as `stopChan <- os.Interrupt`, or a user press
	// CTRL-C in terminal, or others unexpected cases (such as panics).
	verbose("waiting at <-done.")
	if s.msg != "" {
		println(s.msg)
	}
	<-done
	verbose("ended.")
}

func verbose(msg string, args ...any) { //nolint:unparam //no matter
	if VerboseFn != nil {
		VerboseFn(msg, args...)
	}
}

var VerboseFn func(msg string, args ...any)
