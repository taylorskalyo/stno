package cmd

import (
	"fmt"
	"os"

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
		dir, err := stnoDir(notebook)
		if err != nil {
			fmt.Printf("Failed to determine notebook directory: %s.", err.Error())
			os.Exit(1)
		}
		ds, err := datastore.CreateFileStore(dir)
		if err != nil {
			fmt.Printf("Failed to create notebook data store: %s.", err.Error())
			os.Exit(1)
		}

		uuids, err := ds.List()
		if err != nil {
			fmt.Printf("Failed to list notebook entries: %s.", err.Error())
			os.Exit(1)
		}
		tree, _ := toml.Load("")
		for _, uuid := range uuids {
			rc, err := ds.NewReadCloser(uuid)
			if err != nil {
				fmt.Printf("Failed to read notebook entry %s: %s.", uuid, err.Error())
				os.Exit(1)
			}
			defer rc.Close()
			t, err := toml.LoadReader(rc)
			if err != nil {
				fmt.Printf("Invalid toml in entry %s: %s.", uuid, err.Error())
				continue
			}
			tree.Set(uuid, "", false, t)
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
