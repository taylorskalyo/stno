package datastore

import (
	"io"
)

// DataStore represents a persistent data store for notebooks.
type DataStore interface {
	ListEntries(string) ([]string, error)
	NewEntryWriteCloser(string) (io.WriteCloser, error)
	NewEntryReadCloser(string) (io.ReadCloser, error)
	RemoveEntry(string) error
	MoveEntry(string, string) error
}
