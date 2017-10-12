package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/taylorskalyo/stno/datastore"
)

var addCmd = &cobra.Command{
	Use:   "add path",
	Short: "Add an entry to your stno notebook",
	Long: `Add an entry to your stno notebook

New entries will be opened in your default editor (specified by the EDITOR
environment variable). If the entry is not valid TOML, you will be prompted to
revise the entry. This command takes a path. This path is relative to the stno
directory (defaults to ~/.stno). `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add requires a path")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Open entry for writing
		dir, err := homedir.Expand(stnoDir)
		if err != nil {
			fmt.Printf("Could not expand path %s: %s.\n", stnoDir, err.Error())
			os.Exit(1)
		}
		ds := datastore.FileStore{Dir: dir}
		entry, err := ds.NewEntryWriteCloser(args[0])
		if err != nil {
			fmt.Printf("Failed to create entry %s: %s.\n", args[0], err.Error())
			os.Exit(1)
		}

		// Edit entry contents in temporary location
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

		// Write contents to entry
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
