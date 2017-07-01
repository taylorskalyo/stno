package hj

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pelletier/go-toml"
)

// Jrnl represents a journal with entries.
type Jrnl struct {
	Entries []*Entry
	Path    string
}

// Entry represents a journal entry.
type Entry struct {
	toml.Tree
	f *os.File
}

// NewEntry initializes a Entry object.
func (j Jrnl) NewEntry() (e Entry, err error) {
	e.f, err = ioutil.TempFile(j.Path, "")
	if err != nil {
		return
	}
	tree, err := toml.TreeFromMap(make(map[string]interface{}))
	e = Entry{Tree: *tree}
	if err != nil {
		return
	}
	j.Entries = append(j.Entries, &e)
	return
}

// Load opens a journal for reading.
func (j *Jrnl) Load() (err error) {
	listing, err := ioutil.ReadDir(j.Path)
	if err != nil {
		return
	}
	for _, fi := range listing {
		if fi.IsDir() {
			continue
		}
		path := path.Join(j.Path, fi.Name())
		file, err := os.Open(path)
		if err != nil {
			break
		}
		tree, err := toml.LoadFile(path)
		if err != nil {
			break
		}
		e := Entry{
			f:    file,
			Tree: *tree,
		}
		j.Entries = append(j.Entries, &e)
	}
	return
}

// List journal entries.
func (j Jrnl) List() {
	for _, e := range j.Entries {
		fmt.Println(e.f.Name())
		fmt.Println(e.String())
	}
}

// Save saves the contents of an Entry to disk.
func (e Entry) Save() (err error) {
	_, err = e.f.WriteString(e.String())
	return
}

// ApplyTemplate prepopulates an entry with pre defined content.
func (e Entry) ApplyTemplate() {
	e.t.Set("time", time.Now())
}
