package notebook

import (
	"bytes"
	"regexp"
	"text/template"
	"time"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/datastore"
)

const entryIDTemplateStr string = `{{.datetime.Format "2006 Jan 2 15:04:05 MST"}}-{{.title}}`
const entryTemplateStr string = `title    = ""
datetime = {{.datetime}}
notes    = """
"""`

// newTemplateData holds values that can be substituted into a template.
func newEntryTemplateData() map[string]string {
	return map[string]string{
		"datetime": time.Now().Format(time.RFC3339),
	}
}

// Notebook wraps the underlying data store of a stno notebook and provides
// methods for manipulating the contents.
type Notebook struct {
	datastore.DataStore
	entryIDTemplate *template.Template
	entryTemplate   *template.Template
}

// NewEntry generates a new entry based on the notebook's entry template.
func (n Notebook) NewEntry() (*toml.Tree, error) {
	var t *template.Template
	buf := bytes.NewBufferString("")

	// Determine whether to use custom template or default
	if n.entryTemplate != nil {
		t = n.entryTemplate
	} else {
		t = template.Must(template.New("entry").Parse(entryTemplateStr))
	}

	if err := t.Execute(buf, newEntryTemplateData()); err != nil {
		return nil, err
	}
	return toml.LoadReader(buf)
}

// SetEntryTemplate sets a custom entry template. This template allows
// customizing the content used to populate new entries in the notebook.  More
// information about the templating engine used can be found at
// https://golang.org/pkg/text/template/.
func (n *Notebook) SetEntryTemplate(s string) error {
	t, err := template.New("entry").Parse(s)
	if err != nil {
		return err
	}
	n.entryTemplate = t
	return nil
}

// SetEntryIDTemplate sets a custom entry ID template. This template allows
// customizing the ID used to uniquely identify entries within a notebook.
// More information about the templating engine used can be found at
// https://golang.org/pkg/text/template/.
func (n *Notebook) SetEntryIDTemplate(s string) error {
	t, err := template.New("entryID").Parse(s)
	if err != nil {
		return err
	}
	n.entryIDTemplate = t
	return nil
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
