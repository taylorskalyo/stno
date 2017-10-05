package datastore

import (
	"io"
)

// DataStore represents a persistent data store for notebooks.
type DataStore interface {
	List() ([]string, error)
	NewUniqueWriteCloser(string) (string, io.WriteCloser, error)
	NewWriteCloser(string) (io.WriteCloser, error)
	NewReadCloser(string) (io.ReadCloser, error)
}
