package action

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/stno/datastore"
	cli "gopkg.in/urfave/cli.v1"
)

const defaultTemplate string = `title = ""
datetime = {{.DateTime}}
notes = ""`

type templateData struct {
	DateTime string
}

// Add a new notebook entry.
func Add(c *cli.Context) error {
	// Create temporary file
	tmpfile, err := datastore.TempFile("", "stno", ".toml")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	// Write template to file
	t, err := template.New("default").Parse(defaultTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(tmpfile, newTemplateData())
	if err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	// Open file in editor
	openEditor(tmpfile.Name())

	// Lint file
	rc, err := os.Open(tmpfile.Name())
	if err != nil {
		return err
	}
	defer rc.Close()
	tree, err := toml.LoadReader(rc)
	if err != nil {
		return err
	}
	rc.Seek(0, 0)

	// Copy contents from temp file to entry file
	ds, err := datastore.CreateFileStore("/tmp/stno")
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	datetime, ok := tree.Get("datetime").(time.Time)
	if ok {
		buf.WriteString(fmt.Sprintf("%d", datetime.Unix()))
		buf.WriteString("-")
	}
	title, ok := tree.Get("title").(string)
	if ok {
		r := regexp.MustCompile("[^A-Za-z0-9_-]+")
		buf.WriteString(r.ReplaceAllString(title, "-"))
		buf.WriteString("-")
	}
	_, wc, err := ds.NewUniqueWriteCloser(buf.String())
	if err != nil {
		return err
	}
	defer wc.Close()
	_, err = io.Copy(wc, rc)
	if err != nil {
		return err
	}

	return nil
}

// Query a notebook for a list of entries.
func Query(c *cli.Context) error {
	ds, err := datastore.CreateFileStore("/tmp/stno")
	if err != nil {
		return err
	}

	uuids, err := ds.List()
	if err != nil {
		return err
	}
	for i, uuid := range uuids {
		if i != 0 {
			fmt.Println()
		}
		fmt.Println(uuid)
		rc, err := ds.NewReadCloser(uuid)
		if err != nil {
			return err
		}
		io.Copy(os.Stdout, rc)
		rc.Close()
	}

	return nil
}

func openEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "editor"
	}

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}

	editor = args[0]
	args = append(args[1:], path)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// newTemplateData holds values that can be substituted into a template.
func newTemplateData() templateData {
	return templateData{
		DateTime: time.Now().Format(time.RFC3339),
	}
}
