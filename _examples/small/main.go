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

	println(is.InTesting())
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
		WithOnSignalCaught(func(sig os.Signal, wg *sync.WaitGroup) {
			println()
			slog.Info("signal caught", "sig", sig)
			cancel()  // cancel user's loop, see Wait(...)
			wg.Done() // cancel catcher itself
		}).
		Wait(func(stopChan chan<- os.Signal, wgShutdown *sync.WaitGroup) {
			slog.Debug("entering looper's loop...")
			<-ctx.Done()
		})
}
