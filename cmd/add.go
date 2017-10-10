package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	shellquote "github.com/kballard/go-shellquote"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/taylorskalyo/stno/datastore"
	"github.com/taylorskalyo/stno/notebook"
)

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

func stnoDir(name string) (string, error) {
	return homedir.Expand(path.Join("~/.stno", name))
}

// addCmd adds a new entry to the notebook.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a stno entry",
	Long:  `New entries will be opened in your default editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create temporary file
		tmpfile, err := datastore.TempFile("", "stno", ".toml")
		if err != nil {
			fmt.Printf("Failed to create temporary file: %s.\n", err.Error())
			os.Exit(1)
		}
		defer os.Remove(tmpfile.Name())

		// Create notebook
		dir, err := stnoDir(notebookName)
		if err != nil {
			fmt.Printf("Failed to determine notebook directory: %s.\n", err.Error())
			os.Exit(1)
		}
		ds, err := datastore.CreateFileStore(dir)
		if err != nil {
			fmt.Printf("Failed to create notebook data store: %s.\n", err.Error())
			os.Exit(1)
		}
		n := notebook.Notebook{DataStore: ds}

		// Write template to file
		entry, err := n.NewEntry()
		if err != nil {
			fmt.Printf("Failed to apply notebook template: %s.\n", err.Error())
			os.Exit(1)
		}
		_, err = tmpfile.WriteString(entry.String())
		if err != nil {
			fmt.Printf("Failed to write to new notebook entry: %s.\n", err.Error())
			os.Exit(1)
		}
		fi, err := tmpfile.Stat()
		if err != nil {
			fmt.Printf("Failed to stat temporary file %s: %s.\n", tmpfile.Name(), err.Error())
			os.Exit(1)
		}
		oldModTime := fi.ModTime()
		if err := tmpfile.Close(); err != nil {
			fmt.Printf("Failed to close temporary file %s: %s.\n", tmpfile.Name(), err.Error())
			os.Exit(1)
		}

		// Open file in editor
		if err = openEditor(tmpfile.Name()); err != nil {
			fmt.Printf("Failed to open editor: %s.\n", err.Error())
			os.Exit(1)
		}

		rc, err := os.Open(tmpfile.Name())
		if err != nil {
			fmt.Printf("Failed to open temporary file %s: %s.\n", tmpfile.Name(), err.Error())
			os.Exit(1)
		}
		defer rc.Close()

		// Return if there were no changes
		fi, err = rc.Stat()
		if err != nil {
			fmt.Printf("Failed to stat temporary file %s: %s.\n", tmpfile.Name(), err.Error())
			os.Exit(1)
		}
		if oldModTime == fi.ModTime() {
			fmt.Println("Aborting due to empty entry.")
			os.Exit(1)
		}

		// Lint file
		err = entry.LoadReader(rc)
		if err != nil {
			fmt.Printf("Invalid toml: %s.\n", err.Error())
			os.Exit(1)
		}
		rc.Seek(0, 0)

		if err = entry.Save(); err != nil {
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
