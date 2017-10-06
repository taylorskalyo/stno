package notebook

import (
	"text/template"
	"time"
)

const entryIDTemplateStr string = `{{.datetime.Format "2006 Jan 2 15:04:05 MST"}}-{{.title}}`
const entryTemplateStr string = `title = ""
datetime = {{.datetime}}
notes = ""`

type entryTemplateDataFn func() map[string]interface{}

// entryTemplateData holds values that can be substituted into a template.
func entryTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"datetime": time.Now().Format(time.RFC3339),
	}
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
