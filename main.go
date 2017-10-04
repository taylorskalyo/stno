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
			Name:  "notebook, n",
			Usage: "use the notebook named `NAME`",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "load configuration from `FILE`",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "add",
			Usage:  "Adds a notebook entry",
			Action: action.Add,
		},
		{
			Name:   "query",
			Usage:  "Queries notebook entries",
			Action: action.Query,
		},
	}

	app.Run(os.Args)
}
