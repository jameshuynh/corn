package cmd

import (
	"fmt"
	"os/exec"
	"os/user"

	"github.com/jameshuynh/corn/helpers"

	"github.com/spf13/cobra"
)

func migrate(args []string) {
	usr, _ := user.Current()
	dbConfig, adapter := helpers.GenerateDBConfigString()
	if adapter == "psql" {
		adapter = "postgres"
	} else if adapter == "mysql" {
		adapter = "mysql"
	}
	output, _ := exec.Command(
		usr.HomeDir+"/go/bin/mig", "up", adapter, dbConfig, "-d", "db/migrations",
	).CombinedOutput()
	fmt.Println(string(output))
}

var migrateCmd = &cobra.Command{
	Use:   "db:migrate",
	Short: "Generate Different Types",
	Long:  "Generate Different Types",
	Run: func(cmd *cobra.Command, args []string) {
		migrate(args)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
