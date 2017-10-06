package notebook

import (
	"bytes"
	"text/template"

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
func (n Notebook) NewEntry() (Entry, error) {
	var t *template.Template
	var dataFn entryTemplateDataFn
	var entry Entry
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
		return entry, err
	}
	entry.notebook = &n
	if err := entry.LoadReader(buf); err != nil {
		return entry, err
	}
	return entry, nil
}
