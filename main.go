package main

import (
	"os"

	"github.com/taylorskalyo/hj/action"

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
			Name:   "add",
			Usage:  "add a journal entry",
			Action: action.Add,
		},
		{
			Name:   "query",
			Usage:  "filter and display journal entries",
			Action: action.Query,
		},
	}

	app.Run(os.Args)
}
