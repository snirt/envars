package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

)

// var envars *KeePassHandler

// var rootCmd = &cobra.Command{
// 	Use: "envars",
// 	Short: "manage your environment variables",
// 	Long: "This CLI tool help you to keep your environment variable safe and encrypted on your project's directory",
// }

// var addVariables = &cobra.Command{
// 	Use: "add",
// 	Short: "add and modify variables",
// 	Long: "This command adds variables to the database, if the variable already exists, it can be modified",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		envars.AddVariables()
// 	},
// }

// func init() {
// 	envars = New()
// 	rootCmd.AddCommand(addVariables)
// }

func main_() {
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