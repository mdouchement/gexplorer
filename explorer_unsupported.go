//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package gexplorer

type explorer struct{}

func newExplorer(_ RunHandler) *explorer {
	return new(explorer)
}

func (e *Explorer) setView(_ uintptr) {}

func (e *Explorer) importFile(_ ...string) (string, error) {
	return nil, ErrNotAvailable
}

func (e *Explorer) importFiles(_ ...string) ([]string, error) {
	return nil, ErrNotAvailable
}

func (e *Explorer) exportFile(_ string) (string, error) {
	return nil, ErrNotAvailable
}
