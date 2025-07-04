//go:build plan9
// +build plan9

package is

import (
	"os"

	"github.com/hedzr/is/basics"
)

// Raise raises a signal to current process.
//
// It's fairly enough safe and is a better choice versus RaiseSignal.
//
// The common pattern to handle system signals is:
//
//	var stopChan = make(chan os.Signal, 2)
//	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
//	<-stopChan // wait for SIGINT
//
//	// at somewhere you raise it manually
//	stopChan <- syscall.SYSINT
//
// To raise an OS signal is not support on Windows.
func (s signalS) Raise(sig os.Signal) error {
	// TODO cannot work on compiling for plan9
	return basics.Raise(sig)
}
