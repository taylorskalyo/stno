package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/taylorskalyo/stno/datastore"
)

// queryCmd queries a notebook for a list of entries.
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query entries and display the results",
	Long:  `By default query lists all entries.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := homedir.Expand(stnoDir)
		if err != nil {
			fmt.Printf("Could not expand path %s: %s.\n", stnoDir, err.Error())
			os.Exit(1)
		}
		n := datastore.FileStore{Dir: dir}

		uids, err := n.ListEntries("")
		if err != nil {
			fmt.Printf("Failed to list notebook entries: %s.\n", err.Error())
			os.Exit(1)
		}
		tree, _ := toml.Load("")
		for _, uid := range uids {
			rc, err := n.NewEntryReadCloser(uid)
			if err != nil {
				fmt.Printf("Failed to read notebook entry %s: %s.\n", uid, err.Error())
				os.Exit(1)
			}
			defer rc.Close()
			t, err := toml.LoadReader(rc)
			if err != nil {
				fmt.Printf("Invalid toml in entry %s: %s.\n", uid, err.Error())
				continue
			}
			tree.Set(uid, "", false, t)
		}
		fmt.Println(tree.String())
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
