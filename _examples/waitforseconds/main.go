package main

import (
	"context"
	"time"

	"github.com/hedzr/is"
	"github.com/hedzr/is/timing"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := timing.New()
	defer p.CalcNow()

	go func() {
		time.Sleep(3 * time.Second)
		cancel() // stop after 3s instead of waiting for 6s later.
	}()

	is.SignalsEnh().WaitForSeconds(ctx, cancel, 6*time.Second,
		// is.WithCatcherCloser(cancel),
		is.WithCatcherMsg("Press CTRL-C to quit, or waiting for 6s..."),
	)
}
