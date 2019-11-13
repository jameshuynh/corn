package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/jameshuynh/corn/helpers"
	"github.com/spf13/cobra"
)

var sqlboilerCmd = &cobra.Command{
	Use:   "sqlboiler",
	Short: "Generate SQLBoiler",
	Long:  "Generate SQLBoiler",
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := helpers.LoadDBConfig("./config/sqlboiler.toml")
		dbConfig := config.Test
		data := []byte(
			fmt.Sprintf(
				`[psql]
	dbName = "%s"
	host = "%s"
	port = "%d"
	user = "%s"
	password = "%s"
	sslMode = "%s"`,
				dbConfig.Dbname, dbConfig.Host, dbConfig.Port,
				dbConfig.User, dbConfig.Password, dbConfig.Sslmode,
			),
		)
		err := ioutil.WriteFile("config/sqlboiler_test.toml", data, 0755)
		exec.Command("chmod", "-R", "0755", "config/sqlboiler_test.toml").CombinedOutput()
		exec.Command("go", "generate").CombinedOutput()
		fmt.Println(err)
	},
}

func init() {
	rootCmd.AddCommand(sqlboilerCmd)
}
