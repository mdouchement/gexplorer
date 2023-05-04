// Package gexplorer is based on https://github.com/gioui/gio-x/blob/main/explorer
// This version only wants to work with filename/filepath which is more flexible than any other types.
//
//	Given the filename/filepath, the caller of this package is able to do whatever it want.
package gexplorer

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	// ErrUserDecline is returned when the user doesn't select the file.
	ErrUserDecline = errors.New("user exited the file selector without selecting a file")

	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")
)

// RunHandler allows to run a function in the context of another thread.
// Mainly used for https://pkg.go.dev/gioui.org@v0.0.0-20230502183330-59695984e53c/app#Window.Run
type RunHandler func(func())

type result struct {
	filenames []string
	error     error
}

// Explorer facilitates opening OS-native dialogs to choose files and create files.
type Explorer struct {
	id    int32
	mutex sync.Mutex

	// explorer holds OS-Specific content, it varies for each OS.
	*explorer
}

// active holds all explorer currently active, that may necessary for callback functions.
//
// Some OSes (Android, iOS, macOS) may call Golang exported functions as callback, but we need
// someway to link that callback with the respective explorer, in order to give them a response.
//
// In that case, a construction like `callback(..., id int32)` is used. Then, it's possible to get the explorer
// by lookup the active using the callback id.
//
// To avoid hold dead/unnecessary explorer, the active will be removed using `runtime.SetFinalizer` on the related
// Explorer.
var (
	active  = sync.Map{} // map[int32]*explorer
	counter = new(int32)
)

// NewExplorer creates a new Explorer for the given RunHandler.
// The given RunHandler must be unique and you should call NewExplorer
// once per new RunHandler.
func NewExplorer(run RunHandler) (e *Explorer) {
	e = &Explorer{
		explorer: newExplorer(run),
		id:       atomic.AddInt32(counter, 1),
	}

	active.Store(e.id, e.explorer)
	runtime.SetFinalizer(e, func(e *Explorer) { active.Delete(e.id) })

	return e
}

// SetView sets the view/window to the dialod interface.
func (e *Explorer) SetView(v uintptr) {
	e.setView(v)
}

// ChooseFile shows the file selector, allowing the user to select a single file.
// Optionally, it's possible to define which file extensions is supported to
// be selected (such as `.jpg`, `.png`).
//
// Example: ChooseFile(".jpg", ".png") will only accept the selection of files with
// .jpg or .png extensions.
//
// In most known browsers, when user clicks cancel then this function never returns.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile or CreateFile, can happen at the same time, for each Explorer.
func (e *Explorer) ChooseFile(extensions ...string) (string, error) {
	if e == nil {
		return "", ErrNotAvailable
	}

	return e.importFile(extensions...)
}

// ChooseFiles shows the files selector, allowing the user to select multiple files.
// Optionally, it's possible to define which file extensions is supported to
// be selected (such as `.jpg`, `.png`).
//
// Example: ChooseFiles(".jpg", ".png") will only accept the selection of files with
// .jpg or .png extensions.
//
// In most known browsers, when user clicks cancel then this function never returns.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile{,s} or CreateFile, can happen at the same time, for each Explorer.
func (e *Explorer) ChooseFiles(extensions ...string) ([]string, error) {
	if e == nil {
		return nil, ErrNotAvailable
	}

	return e.importFiles(extensions...)
}

// CreateFile opens the file selector, and writes the given content into
// some file, which the use can choose the location.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile or CreateFile, can happen at the same time, for each Explorer.
func (e *Explorer) CreateFile(name string) (string, error) {
	if e == nil {
		return "", ErrNotAvailable
	}

	return e.exportFile(name)
}
