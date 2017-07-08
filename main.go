package main

import (
	"os"

	"github.com/taylorskalyo/hj/journal"
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

	jrnl := journal.Journal{Path: "/tmp/jrnl"}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a journal entry",
			Action: func(c *cli.Context) error {
				entry, err := jrnl.NewEntry()
				if err != nil {
					return err
				}
				entry.ApplyTemplate()
				if err = entry.Save(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "query",
			Usage: "filter and display journal entries",
			Action: func(c *cli.Context) error {
				if err := jrnl.Load(); err != nil {
					return err
				}
				jrnl.List()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
