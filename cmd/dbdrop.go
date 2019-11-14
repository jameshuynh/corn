package cmd

import (
	"fmt"
	"os/exec"

	"github.com/jameshuynh/corn/helpers"
	"github.com/spf13/cobra"
)

func dropDBForEnvironment(environment string) {
	dbConfig := helpers.RetrieveDBConfig(environment)
	if dbConfig.Adapter == "psql" {
		output, _ :=
			exec.Command("dropdb", dbConfig.Dbname).CombinedOutput()
		fmt.Println(string(output))
	} else if dbConfig.Adapter == "mysql" {
		// TODO: handle for mysql
	}
}

func dbDrop(args []string) {
	fmt.Println("drop development db...")
	dropDBForEnvironment("development")
	fmt.Println("drop test db...")
	dropDBForEnvironment("test")
}

var dropDbCmd = &cobra.Command{
	Use:   "db:drop",
	Short: "Generate Different Types",
	Long:  "Generate Different Types",
	Run: func(cmd *cobra.Command, args []string) {
		dbDrop(args)
	},
}

func init() {
	rootCmd.AddCommand(dropDbCmd)
}
