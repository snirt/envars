/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var envars *KeePassHandler
var passwordArg string
var exportArg = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "envars",
	Short: "manage your environment variables",
	Long: `This CLI tool help you to keep your environment variable safe and encrypted on your project's directory`,
	Run: func(cmd *cobra.Command, args []string) { 
		envars = New()
		defer envars.lockDB()
		envars.ListVariables(false);
	},
}

var addCommand = &cobra.Command{
	Use: "add",
	Short: "add and modify variables",
	Long: "This command adds variables to the database, if the variable already exists, it can be modified",
	Run: func(cmd *cobra.Command, args []string) {
		envars = New()
		defer envars.lockDB()
		envars.AddVariables()
	},
}

var removeCommand = &cobra.Command{
	Use: "remove",
	Short: "remove variables from database",
	Long: "This commande removes variables from the database",
	Run: func(cmd *cobra.Command, args []string) {
		envars = New()
		defer envars.lockDB()
		envars.RemoveVariables(args)
	},
}

var listCommand = &cobra.Command{
	Use: "list",
	Short: "List variables",
	Long: `This command list the variables from the database with 'export' prefix. 
		the command 'eval $(envars list -e)' will export the  variables to the session`,
	Run: func(cmd *cobra.Command, args []string) {
		envars = New()
		defer envars.lockDB()
		envars.ListVariables(exportArg)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}


func init() {
	rootCmd.AddCommand(addCommand)
	rootCmd.AddCommand(removeCommand)
	rootCmd.AddCommand(listCommand)

	// rootCmd.PersistentFlags().StringVarP(&passwordArg, "password", "p", "", "envar -p [your_db_password] [COMMAND]")
	listCommand.Flags().BoolVarP(&exportArg, "export", "e", false, "Add export prefix to list")
	
	
	
}


