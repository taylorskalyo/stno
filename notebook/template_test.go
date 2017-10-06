package notebook

import (
	"fmt"
	"reflect"
	"testing"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/notebook"
)

func TestEntryIDTemplate(t *testing.T) {
	testCases := []struct {
		entryContent        string
		customTemplate      string
		expectTemplateError bool
		expectedID          string
	}{
		{ // Default template
			"datetime = 2006-01-02T15:04:05-07:00\ntitle = \"'70s theme party ideas\"",
			"",
			false,
			"2006-Jan-2-15-04-05-0700-70s-theme-party-ideas",
		},
		{ // Custom template
			"foo = \"foo\"\nbar = \"bar\"",
			"{{.foo}}{{.bar}}",
			false,
			"foobar",
		},
		{ // Custom template; missing data
			"bar = \"bar\"\nbaz = \"baz\"",
			"{{.foo}}{{.bar}}",
			false,
			"-no-value-bar",
		},
		{ // Custom template; invalid template
			"foo = \"foo\"",
			"{{foo}}",
			true,
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("template: %s", tc.customTemplate), func(t *testing.T) {
			t.Parallel()
			n := notebook.Notebook{}
			if tc.customTemplate != "" {
				err := n.SetEntryIDTemplate(tc.customTemplate)
				if err != nil && !tc.expectTemplateError {
					t.Fatalf("Unexpected error during SetEntryIDTemplate: %s.", err)
				} else if err != nil && tc.expectTemplateError {
					return
				} else if tc.expectTemplateError {
					t.Fatal("Expected error during SetEntryIDTemplate, but none was thrown.")
				}
			}
			entry, _ := n.NewEntry()
			tree, _ := toml.Load(tc.entryContent)
			entry.Tree = tree
			entryID, err := entry.ID()
			if err != nil {
				t.Fatalf("Unexpected error during EntryID(): %s.", err)
			}
			if entryID != tc.expectedID {
				t.Fatalf("Expected \"%s\", but got \"%s\".", tc.expectedID, entryID)
			}
		})
	}
}

func TestEntryTemplate(t *testing.T) {
	testCases := []struct {
		customTemplate      string
		expectTemplateError bool
		expectTOMLError     bool
		expected            map[string]interface{}
	}{
		{ // Default template
			"",
			false,
			false,
			map[string]interface{}{"title": "", "datetime": "today", "notes": ""},
		},
		{ // Custom template
			"foo = \"foo\"\nbar = \"bar\"",
			false,
			false,
			map[string]interface{}{"foo": "foo", "bar": "bar"},
		},
		{ // Custom template; missing data
			"foo = \"{{.foo}}\"\nbar = \"{{.bar}}\"",
			false,
			false,
			map[string]interface{}{"foo": "<no value>", "bar": "<no value>"},
		},
		{ // Custom template; invalid template
			"foo = \"{{foo}}\"",
			true,
			false,
			map[string]interface{}{},
		},
		{ // Custom template; invalid TOML
			"! invalid TOML",
			false,
			true,
			map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("template: %s", tc.customTemplate), func(t *testing.T) {
			t.Parallel()
			n := notebook.Notebook{
				EntryTemplateDataFn: func() map[string]interface{} {
					return map[string]interface{}{
						"datetime": "\"today\"",
					}
				},
			}
			if tc.customTemplate != "" {
				err := n.SetEntryTemplate(tc.customTemplate)
				if err != nil && !tc.expectTemplateError {
					t.Fatalf("Unexpected error during SetEntryTemplate(): %s.", err)
				} else if err != nil && tc.expectTemplateError {
					return
				} else if tc.expectTemplateError {
					t.Fatal("Expected error during SetEntryTemplate(), but none was thrown.")
				}
			}
			entry, err := n.NewEntry()
			if err != nil && !tc.expectTOMLError {
				t.Fatalf("Unexpected error during NewEntry(): %s.", err)
			} else if err != nil && tc.expectTOMLError {
				return
			} else if tc.expectTOMLError {
				t.Fatal("Expected error during NewEntry(), but none was thrown.")
			}
			actual := entry.ToMap()
			if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("Expected \"%v\", but got \"%v\".", tc.expected, actual)
			}
		})
	}
}
