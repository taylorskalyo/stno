package datastore

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// DataStore represents a persistent data store for Journals.
type DataStore interface {
	List() ([]string, error)
	NewUniqueWriteCloser(string) (string, io.WriteCloser, error)
	NewWriteCloser(string) (io.WriteCloser, error)
	NewReadCloser(string) (io.ReadCloser, error)
}

// FileStore implements DataStore using files and directories.
type FileStore struct {
	Path string
}

// CreateFileStore creates or opens a FileStore at the specified path.
func CreateFileStore(path string) (FileStore, error) {
	var fs FileStore

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fs, err
	}
	fs.Path = path

	return fs, nil
}

// List files in the FileStore.
func (fs FileStore) List() ([]string, error) {
	var entries []string

	err := filepath.Walk(fs.Path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			entries = append(entries, path)
		}
		return nil
	})

	return entries, err
}

// NewUniqueWriteCloser generates a new uuid with an optional prefix and uses it to
// create an io.WriteCloser.
func (fs FileStore) NewUniqueWriteCloser(prefix string) (string, io.WriteCloser, error) {
	f, err := ioutil.TempFile(fs.Path, prefix)
	if err != nil {
		return "", f, err
	}

	return path.Join(fs.Path, f.Name()), f, nil
}

// NewWriteCloser opens or creates the file specified by `path` and returns an
// io.WriteCloser.
func (fs FileStore) NewWriteCloser(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// NewReadCloser opens the file specified by `path` and returns an io.ReadCloser.
func (fs FileStore) NewReadCloser(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
