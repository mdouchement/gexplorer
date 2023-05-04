//go:build darwin && !ios
// +build darwin,!ios

package gexplorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Appkit/AppKit.h>

// Defined on explorer_macos.m file.
extern void exportFile(CFTypeRef viewRef, int32_t id, char * name);
extern void importFile(CFTypeRef viewRef, int32_t id, char * ext);
extern void importFiles(CFTypeRef viewRef, int32_t id, char * ext);
*/
import "C"

import (
	"net/url"
	"strings"
	"unsafe"
)

type explorer struct {
	run    RunHandler
	view   C.CFTypeRef
	result chan result
}

func newExplorer(run RunHandler) *explorer {
	return &explorer{
		run:    run,
		result: make(chan result),
	}
}

func (e *Explorer) setView(v uintptr) {
	e.view = C.CFTypeRef(v)
}

func (e *Explorer) importFile(extensions ...string) (string, error) {
	for i, ext := range extensions {
		extensions[i] = strings.TrimPrefix(ext, ".")
	}

	cextensions := C.CString(strings.Join(extensions, ","))
	e.run(func() {
		C.importFile(e.view, C.int32_t(e.id), cextensions)
	})

	resp := <-e.result
	if resp.error != nil {
		return "", resp.error
	}
	return resp.filenames[0], nil
}

func (e *Explorer) importFiles(extensions ...string) ([]string, error) {
	for i, ext := range extensions {
		extensions[i] = strings.TrimPrefix(ext, ".")
	}

	cextensions := C.CString(strings.Join(extensions, ","))
	e.run(func() {
		C.importFiles(e.view, C.int32_t(e.id), cextensions)
	})

	resp := <-e.result
	if resp.error != nil {
		return nil, resp.error
	}
	return resp.filenames, nil
}

func (e *Explorer) exportFile(name string) (string, error) {
	cname := C.CString(name)
	e.run(func() {
		C.exportFile(e.view, C.int32_t(e.id), cname)
	})

	resp := <-e.result
	if resp.error != nil {
		return "", resp.error
	}
	return resp.filenames[0], nil

}

//export importCallback
func importCallback(id int32, u *C.char) {
	if v, ok := active.Load(id); ok {
		v.(*explorer).result <- newPath([]*C.char{u})
	}
}

//export importsCallback
func importsCallback(id int32, count int32, u **C.char) {
	if v, ok := active.Load(id); ok {
		v.(*explorer).result <- newPath(unsafe.Slice(u, count))
	}
}

//export exportCallback
func exportCallback(id int32, u *C.char) {
	if v, ok := active.Load(id); ok {
		v.(*explorer).result <- newPath([]*C.char{u})
	}
}

func newPath(urls []*C.char) result {
	res := result{
		filenames: make([]string, len(urls)),
	}

	for i, u := range urls {
		name := C.GoString(u)
		if name == "" {
			return result{error: ErrUserDecline}
		}

		uri, err := url.Parse(name)
		if err != nil {
			return result{error: err}
		}

		path, err := url.PathUnescape(uri.Path)
		if err != nil {
			return result{error: err}
		}

		res.filenames[i] = path
	}

	return res
}
