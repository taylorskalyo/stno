package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	toml "github.com/pelletier/go-toml"
	"github.com/taylorskalyo/hj/datastore"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "journal, j",
			Usage: "Load the journal titled `JOURNAL`",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a journal entry",
			Action: func(c *cli.Context) error {
				ds, err := datastore.CreateFileStore("/tmp/jrnl")
				if err != nil {
					return err
				}
				_, w, err := ds.NewUniqueWriter("")
				if err != nil {
					return err
				}
				tpl := `title = ""
datetime = {{.DateTime}}
notes = ""`
				t, err := template.New("defaults").Parse(tpl)
				if err != nil {
					return err
				}
				defaults := struct {
					DateTime string
				}{
					DateTime: time.Now().Format(time.RFC3339),
				}
				err = t.Execute(w, defaults)
				return err
			},
		},
		{
			Name:  "query",
			Usage: "filter and display journal entries",
			Action: func(c *cli.Context) error {
				ds, err := datastore.CreateFileStore("/tmp/jrnl")
				if err != nil {
					return err
				}
				uuids, err := ds.List()
				if err != nil {
					return err
				}
				for _, uuid := range uuids {
					r, err := ds.NewReader(uuid)
					if err != nil {
						return err
					}
					// TODO: preserve formatting; only use toml package for error
					// checking, filtering, and outputting subtrees
					tree, err := toml.LoadReader(r)
					if err != nil {
						return err
					}
					fmt.Println(tree.String())
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
