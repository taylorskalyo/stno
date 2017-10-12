package cmd

import (
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/taylorskalyo/stno/datastore"
)

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
