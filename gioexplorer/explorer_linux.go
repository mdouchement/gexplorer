//go:build linux && !android
// +build linux,!android

package gioexplorer

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/mdouchement/gexplorer"
)

type explorer struct {
	gexplorer *gexplorer.Explorer
}

func newExplorer(_ *app.Window) *explorer {
	return &explorer{
		gexplorer: gexplorer.NewExplorer(nil),
	}
}

func (e *explorer) listenEvents(event event.Event) {
	switch event := event.(type) {
	case app.X11ViewEvent:
		e.gexplorer.SetView(event.Window)
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
