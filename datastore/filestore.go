package datastore

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FileStore implements DataStore using files and directories. Notebooks and
// sections within a notebook are represented by directories. Entries are
// represented by files.
type FileStore struct {
	Dir string
}

// ListEntries returns the UIDs of each entry in the FileStore.
func (fs FileStore) ListEntries(prefix string) ([]string, error) {
	var uids []string

	matches, err := filepath.Glob(path.Join(fs.Dir, prefix+"*"))
	if err != nil {
		return uids, err
	}

	for pathname := range matches {
		fi, err := os.Stat(pathname)
		if err != nil {
			return uids, err
		}
		if fi.IsDir() {
			subUIDS, err := listDir(pathname)
			if err != nil {
				return uids, err
			}
			uids = append(uids, subUIDS)
		} else {
			ext := filepath.Ext(pathname)
			uid := pathname[0 : len(pathname)-len(ext)]
			uids = append(uids, pathname)
		}
	}
	return uids, nil
}

func listDir(dir string) ([]string, error) {
	var uids []string

	err := filepath.Walk(dir, func(pathname string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			ext := filepath.Ext(pathname)
			if strings.ToLower(ext) != ".toml" {
				return nil
			}
			uid := pathname[0 : len(pathname)-len(ext)]
			uids = append(uids, uid)
		}
		return nil
	})
	return uids, err
}

// NewEntryWriteCloser opens or creates an entry's underlying file and returns
// an io.WriteCloser.
func (fs FileStore) NewEntryWriteCloser(uid string) (io.WriteCloser, error) {
	pathname := path.Join(fs.Dir, uid+".toml")
	return os.Create(pathname)
}

// NewEntryReadCloser opens an entry's underlying file and returns an
// io.ReadCloser.
func (fs FileStore) NewEntryReadCloser(uid string) (io.ReadCloser, error) {
	pathname := path.Join(fs.Dir, uid+".toml")
	return os.Open(pathname)
}

// Rename renames (moves) an entry's underlying file from srcUID to destUID
// within the FileStore.
func (fs FileStore) Rename(srcUID, destUID string) error {
	srcPath := path.Join(fs.Dir, srcUID+".toml")
	destPath := path.Join(fs.Dir, destUID+".toml")
	return os.Rename(srcPath, destPath)
}

// Remove removes an entry's underlying file from the FileStore.
func (fs FileStore) Remove(uid string) error {
	pathname := path.Join(fs.Dir, uid+".toml")
	return os.Remove(pathname)
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
