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
	NewUniqueWriter(string) (string, io.Writer, error)
	NewWriter(string) (io.Writer, error)
	NewReader(string) (io.Reader, error)
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

// NewUniqueWriter generates a new uuid with an optional prefix and uses it to
// create an io.Writer.
func (fs FileStore) NewUniqueWriter(prefix string) (string, io.Writer, error) {
	f, err := ioutil.TempFile(fs.Path, prefix)
	if err != nil {
		return "", f, err
	}

	return path.Join(fs.Path, f.Name()), f, nil
}

// NewWriter opens or creates the file specified by `path` and returns an
// io.Writer.
func (fs FileStore) NewWriter(path string) (io.Writer, error) {
	return os.Create(path)
}

// NewReader opens the file specified by `path` and returns an io.Reader.
func (fs FileStore) NewReader(path string) (io.Reader, error) {
	return os.Open(path)
}
