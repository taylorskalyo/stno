package datastore

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// DataStore represents a persistent data store for notebooks.
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

// CreateFileStore creates or opens a FileStore at the specified dir.
func CreateFileStore(dir string) (FileStore, error) {
	var fs FileStore

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fs, err
	}
	fs.Path = dir

	return fs, nil
}

// List UUIDs of entries in the FileStore.
func (fs FileStore) List() ([]string, error) {
	var uuids []string

	err := filepath.Walk(fs.Path, func(pathname string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			base := path.Base(pathname)
			ext := path.Ext(pathname)
			if ext != ".toml" {
				return nil
			}
			uuid := base[0 : len(base)-len(ext)]
			uuids = append(uuids, uuid)
		}
		return nil
	})

	return uuids, err
}

// NewUniqueWriteCloser generates a new uuid with an optional prefix and uses it to
// create an io.WriteCloser.
func (fs FileStore) NewUniqueWriteCloser(prefix string) (string, io.WriteCloser, error) {
	f, err := TempFile(fs.Path, prefix, ".toml")
	if err != nil {
		return "", f, err
	}

	return path.Join(fs.Path, f.Name()), f, nil
}

// NewWriteCloser opens or creates the file with the specified uuid and returns
// an io.WriteCloser.
func (fs FileStore) NewWriteCloser(uuid string) (io.WriteCloser, error) {
	pathname := path.Join(fs.Path, uuid+".toml")
	return os.Create(pathname)
}

// NewReadCloser opens the file with the specified uuid and returns an
// io.ReadCloser.
func (fs FileStore) NewReadCloser(uuid string) (io.ReadCloser, error) {
	pathname := path.Join(fs.Path, uuid+".toml")
	return os.Open(pathname)
}

// TempFile creates a new temporary file in the directory dir with a name
// beginning with prefix and ending in suffix, opens the file for reading and
// writing, and returns the resulting *os.File. TempFile is a modified version
// of ioutil's TempFile that supports suffixes and uses a naive increment
// strategy, instead of random numbers, for uniqueness.
func TempFile(dir, prefix string, suffix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}
	for i := 0; i < 10000; i++ {
		name := path.Join(dir, fmt.Sprintf("%s%d%s", prefix, i, suffix))
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			continue
		}
		break
	}
	return
}
