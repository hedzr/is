package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/hedzr/is"
	"github.com/hedzr/is/basics"
	"github.com/hedzr/is/term/color"
)

func main() {
	defer basics.Close()

	is.RegisterStateGetter("custom", func() bool { return is.InVscodeTerminal() })

	println(is.InTesting())
	println(is.State("in-testing"))
	println(is.State("custom")) // detects a state with custom detector
	println(is.Env().GetDebugLevel())
	if is.InDebugMode() {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})))
	}

	fmt.Printf("%v", color.GetCPT().Translate(`<code>code</code> | <kbd>CTRL</kbd>
		<b>bold / strong / em</b>
		<i>italic / cite</i>
		<u>underline</u>
		<mark>inverse mark</mark>
		<del>strike / del </del>
		<font color="green">green text</font>
`, color.FgDefault))

	ctx, cancel := context.WithCancel(context.Background())
	catcher := is.Signals().Catch()
	catcher.
		WithPrompt("Press CTRL-C to quit...").
		WithOnLoop(dbStarter, cacheStarter, mqStarter).
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			slog.Info("signal caught", "sig", sig)
			cancel() // cancel user's loop, see Wait(...)
		}).
		Wait(func(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
			slog.Debug("entering looper's loop...")
			<-ctx.Done()  // waiting until any os signal caught
			wgDone.Done() // and complete myself
		})
}

func dbStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing database connections...
	// ...
	wgDone.Done()
}

func cacheStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing redis cache connections...
	// ...
	wgDone.Done()
}

func mqStarter(stopChan chan<- os.Signal, wgDone *sync.WaitGroup) {
	// initializing message queue connections...
	// ...
	wgDone.Done()
}
