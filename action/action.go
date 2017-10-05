package action

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	homedir "github.com/mitchellh/go-homedir"
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
		return cli.NewExitError(fmt.Sprintf("Failed to create temporary file: %s.", err.Error()), 1)
	}
	defer os.Remove(tmpfile.Name())

	// Write template to file
	t, err := template.New("default").Parse(defaultTemplate)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to generate notebook template: %s.", err.Error()), 1)
	}
	err = t.Execute(tmpfile, newTemplateData())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to apply notebook template: %s.", err.Error()), 1)
	}
	fi, err := tmpfile.Stat()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to stat temporary file %s: %s.", tmpfile.Name(), err.Error()), 1)
	}
	oldModTime := fi.ModTime()
	if err := tmpfile.Close(); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to close temporary file %s: %s.", tmpfile.Name(), err.Error()), 1)
	}

	// Open file in editor
	if err = openEditor(tmpfile.Name()); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open editor: %s.", err.Error()), 1)
	}

	rc, err := os.Open(tmpfile.Name())
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to open temporary file %s: %s.", tmpfile.Name(), err.Error()), 1)
	}
	defer rc.Close()

	// Return if there were no changes
	fi, err = rc.Stat()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to stat temporary file %s: %s.", tmpfile.Name(), err.Error()), 1)
	}
	if oldModTime == fi.ModTime() {
		return cli.NewExitError("Aborting due to empty entry.", 1)
	}

	// Lint file
	tree, err := toml.LoadReader(rc)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Invalid toml: %s.", err.Error()), 1)
	}
	rc.Seek(0, 0)

	// Copy contents from temp file to entry file
	dir, err := stnoDir(c.GlobalString("notebook"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to determine notebook directory: %s.", err.Error()), 1)
	}
	ds, err := datastore.CreateFileStore(dir)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create notebook data store: %s.", err.Error()), 1)
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
		return cli.NewExitError(fmt.Sprintf("Failed to create notebook entry: %s.", err.Error()), 1)
	}
	defer wc.Close()
	_, err = io.Copy(wc, rc)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to save notebook entry: %s.", err.Error()), 1)
	}

	return nil
}

// Query a notebook for a list of entries.
func Query(c *cli.Context) error {
	dir, err := stnoDir(c.GlobalString("notebook"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to determine notebook directory: %s.", err.Error()), 1)
	}
	ds, err := datastore.CreateFileStore(dir)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create notebook data store: %s.", err.Error()), 1)
	}

	uuids, err := ds.List()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to list notebook entries: %s.", err.Error()), 1)
	}
	for i, uuid := range uuids {
		if i != 0 {
			fmt.Println()
		}
		fmt.Println(uuid)
		rc, err := ds.NewReadCloser(uuid)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to read notebook entry %s: %s.", uuid, err.Error()), 1)
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

func stnoDir(name string) (string, error) {
	return homedir.Expand(path.Join("~/.stno", name))
}
