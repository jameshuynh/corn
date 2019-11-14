package cmd

import (
	"fmt"
	"os/exec"
	"os/user"

	"github.com/jameshuynh/corn/helpers"

	"github.com/spf13/cobra"
)

func roolbackForEnvironment(environment string) {
	usr, _ := user.Current()
	dbConfig, adapter := helpers.GenerateDBConfigString(environment)
	if adapter == "psql" {
		adapter = "postgres"
	} else if adapter == "mysql" {
		adapter = "mysql"
	}
	output, _ := exec.Command(
		usr.HomeDir+"/go/bin/mig", "down", adapter, dbConfig, "-d", "db/migrations",
	).CombinedOutput()
	fmt.Println(string(output))
}

func rollback(args []string) {
	roolbackForEnvironment("development")
	roolbackForEnvironment("test")
}

var rollbackCmd = &cobra.Command{
	Use:   "db:rollback",
	Short: "Generate Different Types",
	Long:  "Generate Different Types",
	Run: func(cmd *cobra.Command, args []string) {
		rollback(args)
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
