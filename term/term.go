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
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"sync"

	"github.com/hedzr/is/basics"
	"github.com/hedzr/is/dir"
	"github.com/hedzr/is/dirs"
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

func IsENOTTY(err error) bool { return errIsENOTTY(err) }

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

type PromptModeConfig struct {
	Name              string // used for binding history records
	WelcomeText       string
	PromptText        string
	ReplyText         string
	MainLooperHandler LooperFunc
	PostInitTerminal  func(t *term.Terminal)
	MaxHistoryEntries int
}

func MakeNewTerm(ctx context.Context, config *PromptModeConfig) (deferFunc func(), err error) {
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	term := term.NewTerminal(screen, config.PromptText)
	term.SetPrompt(string(term.Escape.Red) + config.PromptText + string(term.Escape.Reset))

	rePrefix := string(term.Escape.Cyan) + config.ReplyText + string(term.Escape.Reset)

	exitChan := make(chan struct{}, 3)
	deferFunc = func() { close(exitChan) }

	if config.MaxHistoryEntries <= 1 {
		config.MaxHistoryEntries = 1000
	}
	if config.MaxHistoryEntries > 12000 {
		config.MaxHistoryEntries = 12000
	}
	if historyName := config.Name; historyName != "" {
		historyDir := dirs.DataDir(config.Name, "is.term.history")
		if err = dir.EnsureDir(historyDir); err != nil {
			return
		}
		historyFile := path.Join(historyDir, "history.list")
		if dir.FileExists(historyFile) {
			var file *os.File
			file, err = os.Open(historyFile)
			if err != nil {
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				// fmt.Println(scanner.Text())
				term.History.Add(line)
			}

			if err = scanner.Err(); err != nil {
				return
			}

			// fmt.Printf("%d entries of %q loaded.\n", term.History.Len(), historyFile)
			slog.Debug("history file has been loaded.\n", "enteries", term.History.Len(), "file", historyFile)
		}
		defer func() {
			var file *os.File
			file, err = os.Create(historyFile)
			if err != nil {
				return
			}
			defer file.Close()
			l := term.History.Len()
			for i := max(l-config.MaxHistoryEntries, 0); i < l; i++ {
				line := term.History.At(i)
				if _, err = file.WriteString(line); err != nil {
					return
				}
				file.WriteString("\n")
			}
			// slog.Debug("%d entries of %q written.\n", term.History.Len(), historyFile)
			slog.Debug("history file has been written.\n", "enteries", term.History.Len(), "file", historyFile)
		}()
	}

	if fn := config.PostInitTerminal; fn != nil {
		fn(term)
	}

	catcher := basics.Catch()
	catcher.
		// WithVerboseFn(func(msg string, args ...any) {
		// 	// logz.WithSkip(2).PrintlnContext(ctx, fmt.Sprintf("[verbose] %s\n", fmt.Sprintf(msg, args...)))
		// }).
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			// logz.Debug("signal caught", "sig", sig)
			exitChan <- struct{}{}
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			if config.WelcomeText != "" {
				_, _ = fmt.Fprintln(term, config.WelcomeText)
			}
			if fn := config.MainLooperHandler; fn != nil {
				err = fn(ctx, term, rePrefix, exitChan, closer)
			}
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
