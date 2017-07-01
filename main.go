package main

import (
	"os"

	"github.com/taylorskalyo/hj/hj"
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

	jrnl := hj.Jrnl{Path: "/tmp/jrnl"}

	app.Commands = []cli.Command{
		{
			Name:  "add",
			Usage: "add a journal entry",
			Action: func(c *cli.Context) error {
				e, err := jrnl.NewEntry()
				if err != nil {
					return err
				}
				e.ApplyTemplate()
				if err = e.Save(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "query",
			Usage: "filter and display journal entries",
			Action: func(c *cli.Context) error {
				if err := jrnl.Open(); err != nil {
					return err
				}
				jrnl.List()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
