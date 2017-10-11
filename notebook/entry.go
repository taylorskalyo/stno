package notebook

import (
	"io"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/notebook"
)

// Entry wraps the underlying TOML tree of a stno entry and provides methods
// for manipulating the contents.
type Entry struct {
	notebook *notebook.Notebook
	UID      string
}

// NewEntry generates a new entry based on the notebook's entry template.
func (n Notebook) NewEntry(uid string) Entry {
	return Entry{notebook: n, UID: uid}
}

// WriteString persists an entry to the notebook's underlying data store.
func (e *Entry) WriteString(contents string) error {
	wc, err := e.notebook.NewWriteCloser(e.uid)
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err := toml.Load(contents)
	if err != nil {
		return err
	}

	_, err = io.WriteString(wc, contents)
	return err
}
