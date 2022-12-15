package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var envars *KeePassHandler

func init() {
	envars = New()
}

func main() {
	defer envars.lockDB()
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Category: "Variables",
				Name:     "Add new variables",
				Aliases:  []string{"add"},
				Action: func(cCtx *cli.Context) error {
					envars.AddVariables()
					return nil
				},
			},
			{
				Category: "Variables",
				Name:     "Export all variables",
				Aliases:  []string{"export"},
				Action: func(cCtx *cli.Context) error {
					envars.ListVariables(true)
					return nil
				},
			},
            {
				Category: "Variables",
				Name:     "Remove variables",
				Aliases:  []string{"rm"},
				Action: func(cCtx *cli.Context) error {
                    vars := cCtx.Args().Slice()
					envars.RemoveVariables(vars)
					return nil
				},
			},
		},
		Action: func(cCtx *cli.Context) error {
			envars.ListVariables(false)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}