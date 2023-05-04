//go:build darwin && !ios
// +build darwin,!ios

package gioexplorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Appkit/AppKit.h>
*/
import "C"

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/mdouchement/gexplorer"
)

type explorer struct {
	window    *app.Window
	gexplorer *gexplorer.Explorer
}

func newExplorer(w *app.Window) *explorer {
	return &explorer{
		window:    w,
		gexplorer: gexplorer.NewExplorer(w.Run),
	}
}

func (e *explorer) listenEvents(event event.Event) {
	switch event := event.(type) {
	case app.ViewEvent:
		e.gexplorer.SetView(event.View)
	}
}

func (e *explorer) importFile(extensions ...string) (string, error) {
	return e.gexplorer.ChooseFile(extensions...)
}

func (e *explorer) importFiles(extensions ...string) ([]string, error) {
	return e.gexplorer.ChooseFiles(extensions...)
}

func (e *explorer) exportFile(name string) (string, error) {
	return e.gexplorer.CreateFile(name)
}
