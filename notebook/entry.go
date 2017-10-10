package notebook

import (
	"bytes"
	"io"
	"regexp"
	"text/template"

	toml "github.com/pelletier/go-toml"
)

// Entry wraps the underlying TOML tree of a stno entry and provides methods
// for manipulating the contents.
type Entry struct {
	*toml.Tree
	notebook *Notebook
}

// ID generates an ID for the given entry based on the notebook's entry ID
// template (see SetEntryIDTemplate). This ID is not necessarily unique between
// entries, but it will be used to generate a unique identifier.
func (e Entry) ID() (string, error) {
	var t *template.Template
	buf := bytes.NewBufferString("")

	// Determine whether to use custom template or default
	if e.notebook.entryIDTemplate != nil {
		t = e.notebook.entryIDTemplate
	} else {
		t = template.Must(template.New("entryID").Parse(entryIDTemplateStr))
	}

	if err := t.Execute(buf, e.ToMap()); err != nil {
		return "", err
	}

	// Replace any non-alphanumeric characters
	r := regexp.MustCompile("-*[^A-Za-z0-9_-]+-*")
	return r.ReplaceAllString(buf.String(), "-"), nil
}

// LoadReader populates an entry with content from a Reader.
func (e *Entry) LoadReader(r io.Reader) error {
	tree, err := toml.LoadReader(r)
	if err != nil {
		return err
	}
	e.Tree = tree
	return nil
}
