//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package gioexplorer

import (
	"io"

	"gioui.org/app"
	"gioui.org/io/event"
)

type explorer struct{}

func newExplorer(w *app.Window) *explorer {
	return new(explorer)
}

func (e *Explorer) listenEvents(_ event.Event) {}

func (e *Explorer) exportFile(_ string) (string, error) {
	return nil, ErrNotAvailable
}

func (e *Explorer) importFile(_ ...string) (string, error) {
	return nil, ErrNotAvailable
}

func (e *Explorer) importFiles(_ ...string) ([]string, error) {
	return nil, ErrNotAvailable
}
