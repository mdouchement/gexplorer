package gioexplorer

import (
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
)

// Explorer facilitates opening OS-native dialogs to choose files and create files.
type Explorer struct {
	*explorer
}

// NewExplorer creates a new Explorer for the given *app.Window.
// The given app.Window must be unique and you should call NewExplorer
// once per new app.Window.
//
// It's mandatory to use Explorer.ListenEvents on the same *app.Window.
func NewExplorer(w *app.Window) (e *Explorer) {
	return &Explorer{
		explorer: newExplorer(w),
	}
}

// ListenEvents must get all the events from Gio, in order to get the GioView. You must
// include that function where you listen for Gio events.
//
// Similar as:
//
//	select {
//	case e := <-window.Events():
//
//		explorer.ListenEvents(e)
//		switch e := e.(type) {
//			(( ... your code ...  ))
//		}
//	}
func (e *Explorer) ListenEvents(evt event.Event) {
	if e == nil {
		return
	}
	e.listenEvents(evt)
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
// ChooseFile or CreateFile, can happen at the same time, for each app.Window/Explorer.
func (e *Explorer) ChooseFile(extensions ...string) (string, error) {
	return e.importFile(extensions...)
}

// ChooseFileIO shows the file selector, allowing the user to select a single file.
// Optionally, it's possible to define which file extensions is supported to
// be selected (such as `.jpg`, `.png`).
//
// Example: ChooseFileIO(".jpg", ".png") will only accept the selection of files with
// .jpg or .png extensions.
//
// In some platforms the resulting `io.ReadCloser` is a `os.File`, but it's not
// a guarantee.
//
// In most known browsers, when user clicks cancel then this function never returns.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile or CreateFile, can happen at the same time, for each app.Window/Explorer.
func (e *Explorer) ChooseFileIO(extensions ...string) (io.ReadCloser, error) {
	filename, err := e.ChooseFile(extensions...)
	if err != nil {
		return nil, err
	}

	return os.Open(filename)
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
// ChooseFile{,s} or CreateFile, can happen at the same time, for each app.Window/Explorer.
func (e *Explorer) ChooseFiles(extensions ...string) ([]string, error) {
	return e.importFiles(extensions...)
}

// ChooseFilesIO shows the files selector, allowing the user to select multiple files.
// Optionally, it's possible to define which file extensions is supported to
// be selected (such as `.jpg`, `.png`).
//
// Example: ChooseFilesIO(".jpg", ".png") will only accept the selection of files with
// .jpg or .png extensions.
//
// In some platforms the resulting `io.ReadCloser` is a `os.File`, but it's not
// a guarantee.
//
// In most known browsers, when user clicks cancel then this function never returns.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile{,s} or CreateFile, can happen at the same time, for each app.Window/Explorer.
func (e *Explorer) ChooseFilesIO(extensions ...string) ([]io.ReadCloser, error) {
	filenames, err := e.ChooseFiles(extensions...)
	if err != nil {
		return nil, err
	}

	readers := make([]io.ReadCloser, len(filenames))
	for i, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		readers[i] = f
	}

	return readers, nil
}

// CreateFile opens the file selector, and writes the given content into
// some file, which the use can choose the location.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile or CreateFile, can happen at the same time, for each Explorer.
func (e *Explorer) CreateFile(name string) (string, error) {
	return e.exportFile(name)
}

// CreateFileIO opens the file selector, and writes the given content into
// some file, which the use can choose the location.
//
// It's important to close the `io.WriteCloser`. In some platforms the
// file will be saved only when the writer is closer.
//
// In some platforms the resulting `io.WriteCloser` is a `os.File`, but it's not
// a guarantee.
//
// It's a blocking call, you should call it on a separated goroutine. For most OSes, only one
// ChooseFile or CreateFile, can happen at the same time, for each app.Window/Explorer.
func (e *Explorer) CreateFileIO(name string) (io.WriteCloser, error) {
	filename, err := e.CreateFile(name)
	if err != nil {
		return nil, err
	}

	return os.Create(filename)
}
