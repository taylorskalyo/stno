package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/taylorskalyo/stno/datastore"
)

const stnoDir string = "~/.stno"

func editString(str string) (string, error) {
	// Create temporary file
	tmpfile, err := datastore.TempFile("", "stno", ".toml")
	if err != nil {
		return str, err
	}
	defer os.Remove(tmpfile.Name())
	if _, err = io.WriteString(tmpfile, str); err != nil {
		return str, err
	}

	fi, err := tmpfile.Stat()
	if err != nil {
		return str, err
	}
	oldModTime := fi.ModTime()
	if err := tmpfile.Close(); err != nil {
		return str, err
	}
	tmpfile.Close()

	// Open file in editor
	if err = openEditor(tmpfile.Name()); err != nil {
		return str, err
	}

	rc, err := os.Open(tmpfile.Name())
	if err != nil {
		return str, err
	}
	defer rc.Close()

	// Return original string if there were no changes
	fi, err = rc.Stat()
	if err != nil {
		return str, err
	}
	if oldModTime == fi.ModTime() {
		fmt.Println("Aborting due to no changes.")
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rc)
	if err != nil {
		return str, err
	}
	return buf.String(), nil
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

func confirm() bool {
	fmt.Printf("(y/N): ")
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
