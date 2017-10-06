package notebook

import (
	"bytes"
	"regexp"
	"text/template"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/datastore"
)

// Notebook wraps the underlying data store of a stno notebook and provides
// methods for manipulating the contents.
type Notebook struct {
	datastore.DataStore
	entryIDTemplate     *template.Template
	entryTemplate       *template.Template
	EntryTemplateDataFn entryTemplateDataFn
}

// NewEntry generates a new entry based on the notebook's entry template.
func (n Notebook) NewEntry() (*toml.Tree, error) {
	var t *template.Template
	var dataFn entryTemplateDataFn
	buf := bytes.NewBufferString("")

	// Determine whether to use custom template or default
	if n.entryTemplate != nil {
		t = n.entryTemplate
	} else {
		t = template.Must(template.New("entry").Parse(entryTemplateStr))
	}
	if n.EntryTemplateDataFn != nil {
		dataFn = n.EntryTemplateDataFn
	} else {
		dataFn = entryTemplateData
	}

	if err := t.Execute(buf, dataFn()); err != nil {
		return nil, err
	}
	return toml.LoadReader(buf)
}

// EntryID generates an ID for the given entry based on the notebook's entry ID
// template (see SetEntryIDTemplate). This ID is not necessarily unique between
// entries, but it will be used to generate a unique identifier.
func (n Notebook) EntryID(entry *toml.Tree) (string, error) {
	var t *template.Template
	buf := bytes.NewBufferString("")

	// Determine whether to use custom template or default
	if n.entryIDTemplate != nil {
		t = n.entryIDTemplate
	} else {
		t = template.Must(template.New("entryID").Parse(entryIDTemplateStr))
	}

	if err := t.Execute(buf, entry.ToMap()); err != nil {
		return "", err
	}

	// Replace any non-alphanumeric characters
	r := regexp.MustCompile("-*[^A-Za-z0-9_-]+-*")
	return r.ReplaceAllString(buf.String(), "-"), nil
}
