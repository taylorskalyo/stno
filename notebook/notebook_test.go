package notebook

import (
	"reflect"
	"sort"
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

func TestCustomValidEntryIDTemplate(t *testing.T) {
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

func TestCustomMissingEntryIDTemplate(t *testing.T) {
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

func TestCustomInvalidEntryIDTemplate(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `{{foo}}`
	if err := n.SetEntryIDTemplate(customTemplate); err == nil {
		t.Fatal("Expected error, but none was thrown.")
	}
}

func TestDefaultEntry(t *testing.T) {
	n := notebook.Notebook{}
	entry, err := n.NewEntry()
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := []string{"datetime", "notes", "title"}
	actual := entry.Keys()
	sort.Strings(actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected \"%v\", but got \"%v\".", expected, actual)
	}
}

func TestCustomValidEntryTemplate(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `foo = "fighter"
bar = "none"`
	if err := n.SetEntryTemplate(customTemplate); err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	entry, err := n.NewEntry()
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := map[string]string{"foo": "fighter", "bar": "none"}
	actual := entry.ToMap()
	if expected["foo"] != actual["foo"] || expected["bar"] != actual["bar"] {
		t.Fatalf("Expected \"%v\", but got \"%v\".", expected, actual)
	}
}

func TestCustomMissingEntryTemplate(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `foo = "{{.foo}}"
bar = "{{.bar}}"`
	if err := n.SetEntryTemplate(customTemplate); err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	entry, err := n.NewEntry()
	if err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	expected := map[string]string{"foo": "<no value>", "bar": "<no value>"}
	actual := entry.ToMap()
	if expected["foo"] != actual["foo"] || expected["bar"] != actual["bar"] {
		t.Fatalf("Expected \"%v\", but got \"%v\".", expected, actual)
	}
}

func TestCustomInvalidEntryTemplate(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `{{foo}}`
	if err := n.SetEntryTemplate(customTemplate); err == nil {
		t.Fatal("Expected error, but none was thrown.")
	}
}

func TestCustomInvalidTOMLEntryTemplate(t *testing.T) {
	n := notebook.Notebook{}
	customTemplate := `!invalid TOML`
	if err := n.SetEntryTemplate(customTemplate); err != nil {
		t.Fatalf("Unexpected error: %s.", err)
	}
	_, err := n.NewEntry()
	if err == nil {
		t.Fatal("Expected error, but none was thrown.")
	}
}
