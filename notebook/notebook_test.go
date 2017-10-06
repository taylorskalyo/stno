package notebook

import (
	"testing"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/notebook"
)

func TestDefaultEntryID(t *testing.T) {
	n := notebook.Notebook{}
	entry, _ := toml.Load(`datetime = 2006-01-02T15:04:05-07:00
		title = "'70s theme party ideas"`)
	entryID, err := n.EntryID(entry)
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := "2006-Jan-2-15-04-05-0700-70s-theme-party-ideas"
	if entryID != expected {
		t.Fatalf("Expected \"%s\", but got \"%s\".", expected, entryID)
	}
}

func TestCustomValidEntryID(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `{{.foo}}{{.bar}}`
	if err := n.SetEntryIDTemplate(customTemplate); err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	entry, _ := toml.Load(`foo = "foo"
	bar = "bar"`)
	entryID, err := n.EntryID(entry)
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := "foobar"
	if entryID != expected {
		t.Fatalf("Expected \"%s\", but got \"%s\".", expected, entryID)
	}
}

func TestCustomMissingEntryID(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `{{.foo}}{{.bar}}`
	if err := n.SetEntryIDTemplate(customTemplate); err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	entry, _ := toml.Load(`bar = "bar"
	baz = "baz"`)
	entryID, err := n.EntryID(entry)
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := "-no-value-bar"
	if entryID != expected {
		t.Fatalf("Expected \"%s\", but got \"%s\".", expected, entryID)
	}
}

func TestCustomInvalidEntryID(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `{{foo}}`
	if err := n.SetEntryIDTemplate(customTemplate); err == nil {
		t.Fatal("Expected error, but none was thrown.")
	}
}
