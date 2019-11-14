package cmd

import (
	"fmt"
	"os/exec"

	"github.com/jameshuynh/corn/helpers"
	"github.com/spf13/cobra"
)

func createForEnvironment(environment string) {
	dbConfig := helpers.RetrieveDBConfig(environment)
	if dbConfig.Adapter == "psql" {
		output, _ :=
			exec.Command("createdb", dbConfig.Dbname, "--encoding=utf8").CombinedOutput()
		fmt.Println(string(output))
	} else if dbConfig.Adapter == "mysql" {
		// TODO: handle for mysql
	}
}

func dbCreate(args []string) {
	fmt.Println("create development db...")
	createForEnvironment("development")
	fmt.Println("create test db...")
	createForEnvironment("test")
}

var createDbCmd = &cobra.Command{
	Use:   "db:create",
	Short: "Generate Different Types",
	Long:  "Generate Different Types",
	Run: func(cmd *cobra.Command, args []string) {
		dbCreate(args)
	},
}

func init() {
	rootCmd.AddCommand(createDbCmd)
}
