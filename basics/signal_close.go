package basics

import (
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
// For example:
//
//	basics.VerboseFn = t.Logf
//	env.Signals().
//		Catch().
//		WithPrompt().
//		Wait(func(stopChan chan<- os.Signal) {
//			basics.VerboseFn("[cb] raising interrupt after a second...")
//			time.Sleep(2500 * time.Millisecond)
//			stopChan <- os.Interrupt
//			basics.VerboseFn("[cb] raised.")
//		})
//
// A simple details can be found at:
//
//	https://www.developer.com/languages/os-signals-go/
//
// Your logic to shutdown the main loop gracefully could be:
//
//	type appS struct{}
//	func (s *appS) DoClose(stopChan chan<- os.Signal) {
//	     stopChan <- os.Interrupt
//	}
//	var app appS
//
//	env.Signals().
//	    Catch().
//	    Wait(app.DoClose)
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
	// An empty call means the default set are applying:
	//    signals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP}
	// These signals are hard-coded.
	WithSignals(signals ...os.Signal) Catcher
	// WithCloseHandlers gives the extra Peripheral 's a chance with safely shutting down.
	//
	// The better way is registering them with env.Closers().RegisterPeripheral(p). All
	// registered peripherals will be released/closed in the ending of Wait.
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
	WithOnSignalCaught(cb ...OnSignalCaught) Catcher
	WithOnLoop(cb ...OnLooper) Catcher
	WithVerboseFn(cb func(msg string, args ...any)) Catcher
	// Wait get the current thread blocked on reading done chan till a os
	// signal break it.
	//
	// You may send a signal (commonly like os.Interrupt) to stopChan to stop
	// the blocked state programmatically.
	Wait(stopperHandlers ...OnLooper)
}

type OnSignalCaught func(sig os.Signal, wgShutdown *sync.WaitGroup)
type OnLooper func(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup)

type catsig struct {
	signals        []os.Signal
	closers        []Peripheral
	onCaught       []OnSignalCaught
	looperHandlers []OnLooper
	msg            string
}

func (s *catsig) WithCloseHandlers(onFinish ...Peripheral) Catcher {
	s.closers = append(s.closers, onFinish...)
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

func (s *catsig) WithVerboseFn(cb func(msg string, args ...any)) Catcher {
	VerboseFn = cb
	return s
}

func (s *catsig) Wait(looperHandlerS ...OnLooper) {
	s.looperHandlers = append(s.looperHandlers, looperHandlerS...)

	var wgInitialized sync.WaitGroup
	wgInitialized.Add(len(s.looperHandlers))

	var closed int32
	done := make(chan struct{})
	defer func() {
		verbose("closing onFinish...")
		for _, f := range s.closers {
			f.Close()
		}
		verbose("closing registered closers...")
		closers.Close()
		if atomic.CompareAndSwapInt32(&closed, 0, 1) {
			verbose("closing done chan finally...")
			close(done)
		}
	}()

	c := make(chan os.Signal, 8) // allow buffering some signals
	signals := s.signals
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT}
	}
	signal.Notify(c, signals...) //nolint:govet //whyNoLint for why

	var wgForShutdown sync.WaitGroup
	wgForShutdown.Add(len(s.onCaught))

	for _, f := range s.looperHandlers {
		go func(c chan os.Signal, wgInitialized, wgForShutdown *sync.WaitGroup, f OnLooper) {
			wgInitialized.Done()
			f(c, wgForShutdown)
		}(c, &wgInitialized, &wgForShutdown, f)
	}
	wgInitialized.Wait()

	verbose("all looper(s) ran")

	go func(wg *sync.WaitGroup) {
		sig := <-c
		for _, cb := range s.onCaught {
			if cb != nil {
				cb(sig, wg)
			}
		}
		if len(s.onCaught) == 0 {
			println()
			verbose("signal caught", "sig", sig)
		}
		wg.Wait()
		done <- struct{}{}
	}(&wgForShutdown)

	// enter main loop here till someone raise a signal from looperHandlers
	// by triggering such as  `stopChan <- os.Interrupt`, or a user press
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
