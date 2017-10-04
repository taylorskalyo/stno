package main

import (
	"os"

	"github.com/taylorskalyo/stno/action"

	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "journal, j",
			Usage: "load the journal titled `JOURNAL`",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Adds a journal entry",
			Action: action.Add,
		},
		{
			Name:   "query",
			Usage:  "Queries journal entries",
			Action: action.Query,
		},
	}

	app.Run(os.Args)
}
