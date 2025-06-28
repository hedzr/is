// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package term provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
//	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
//	if err != nil {
//	        panic(err)
//	}
//	defer term.Restore(int(os.Stdin.Fd()), oldState)
//
// Note that on non-Unix systems os.Stdin.Fd() may not be 0.
package term

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hedzr/is/basics"
	"golang.org/x/term"
)

// // IsTerminal return true if the file descriptor is terminal.
// func IsTerminal(fd uintptr) bool {
// 	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
// 	return err == nil
// }

// IsCygwinTerminal return true if the file descriptor is a cygwin or msys2
// terminal. This is also always false on this environment.
func IsCygwinTerminal(fd uintptr) bool {
	return false
}

// IsTerminal returns whether the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// MakeRaw puts the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*term.State, error) {
	return term.MakeRaw(fd)
}

func MakeRawWrapped() (deferFunc func(), err error) {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		return func() {}, fmt.Errorf("stdin/stdout should be terminal")
	}

	var oldState *term.State
	oldState, err = term.MakeRaw(0)
	if err != nil {
		if !errIsENOTTY(err) {
			return func() {}, err
		}
		// if !errors.Is(err, syscall.ENOTTY) {
		// 	return func() {}, err
		// }
	}
	deferFunc = func() {
		if e := recover(); e != nil {
			if err == nil {
				if e1, ok := e.(error); ok {
					err = e1
				} else {
					err = fmt.Errorf("%v", e)
				}
			} else {
				err = fmt.Errorf("%v | %v", e, err)
			}
		}

		if e1 := term.Restore(0, oldState); e1 != nil {
			if err == nil {
				err = e1
			} else {
				err = fmt.Errorf("%v | %v", e1, err)
			}
		}
	}
	return
}

type SmallTerm interface {
	io.Writer
	ReadLine() (string, error)
}

type LooperFunc func(ctx context.Context, tty SmallTerm, replyPrefix string, exitChan <-chan struct{}, closer func()) (err error)

func MakeNewTerm(ctx context.Context, welcomeString, promptString, replyPrefix string, looper LooperFunc) (deferFunc func(), err error) {
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := term.NewTerminal(screen, promptString)
	term.SetPrompt(string(term.Escape.Red) + promptString + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + replyPrefix + string(term.Escape.Reset)

	exitChan := make(chan struct{}, 3)
	deferFunc = func() { close(exitChan) }

	catcher := basics.Catch()
	catcher.
		WithVerboseFn(func(msg string, args ...any) {
			// logz.WithSkip(2).PrintlnContext(ctx, fmt.Sprintf("[verbose] %s\n", fmt.Sprintf(msg, args...)))
		}).
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			// logz.Debug("signal caught", "sig", sig)
			exitChan <- struct{}{}
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			if welcomeString != "" {
				_, _ = fmt.Fprintln(term, welcomeString)
			}
			err = looper(ctx, term, rePrefix, exitChan, closer)
		})
	return
}

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*term.State, error) {
	return term.GetState(fd)
}

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, oldState *term.State) error {
	return term.Restore(fd, oldState)
}

// GetSize returns the visible dimensions of the given terminal.
//
// These dimensions don't include any scrollback buffer height.
func GetSize(fd int) (width, height int, err error) {
	return term.GetSize(fd)
}

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) {
	return term.ReadPassword(fd)
}
