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
	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/taylorskalyo/stno/datastore"
)

const stnoDir string = "~/.stno"

// addCmd adds a new entry to the notebook.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a stno entry",
	Long:  `New entries will be opened in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create entry
		dir, err := homedir.Expand(stnoDir)
		if err != nil {
			fmt.Printf("Could not expand path %s: %s.\n", stnoDir, err.Error())
			os.Exit(1)
		}
		n := datastore.FileStore{Dir: dir}
		entry, err := n.NewEntryWriteCloser("entry")
		if err != nil {
			fmt.Printf("Failed to create entry %s: %s.\n", "entry", err.Error())
			os.Exit(1)
		}

		contents := ""
		for {
			contents, err = editString(contents)
			if _, err = toml.Load(contents); err != nil {
				fmt.Printf("Could not parse TOML: %s. Try again? ", err.Error())
				if !confirm() {
					break
				}
			} else {
				break
			}
		}

		if _, err = io.WriteString(entry, contents); err != nil {
			fmt.Printf("Failed to save notebook entry: %s.\n", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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
