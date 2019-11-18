package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/jameshuynh/corn/helpers"
	"github.com/spf13/cobra"
)

func searchAndReplaceSQLBoilerConfig() {
	replacers := make(map[string]string)
	replacers["wd = wd + strings.Repeat(\"/..\", outputDirDepth)"] =
		"wd = wd + strings.Repeat(\"/../config\", outputDirDepth)"
	replacers["\"sqlboiler\"))"] = "\"database\"))"
	replacers["viper.SetConfigName(\"sqlboiler\")"] =
		"viper.SetConfigName(\"database\")"
	replacers["GetString(\"psql."] = "GetString(\"test."
	replacers["GetInt(\"psql."] = "GetInt(\"test."
	replacers["\"psql."] = "\"test."
	helpers.SearchAndReplaceFiles(".", replacers)
}

func boilerGeneratorForMysql() {
	// TODO
}

func boilerGeneratorForPsql() {
	config, _ := helpers.LoadDBConfig("./config/database.toml")
	dbConfig := config.Test
	data := []byte(
		fmt.Sprintf(
			`[psql]
	dbName = "%s"
	host = "%s"
	port = "%d"
	user = "%s"
	password = "%s"
	sslMode = "%s"
	blacklist = ["mig_migrations"]`,
			dbConfig.Dbname, dbConfig.Host, dbConfig.Port,
			dbConfig.User, dbConfig.Password, dbConfig.Sslmode,
		),
	)
	err := ioutil.WriteFile("config/database_test.toml", data, 0755)
	if err != nil {
		panic(err)
	}
	exec.Command(
		"chmod", "-R", "0755", "config/database_test.toml",
	).CombinedOutput()
	exec.Command(
		"sqlboiler", "-c", "config/database_test.toml", "--wipe", "psql",
		"--add-global-variants",
	).CombinedOutput()
	exec.Command("rm", "-rf", "config/database_test.toml").CombinedOutput()

	searchAndReplaceSQLBoilerConfig()
}

var sqlboilerCmd = &cobra.Command{
	Use:   "sqlboiler",
	Short: "Generate SQLBoiler",
	Long:  "Generate SQLBoiler",
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := helpers.LoadDBConfig("./config/database.toml")
		if config.Test.Adapter == "psql" {
			boilerGeneratorForPsql()
		} else if config.Test.Adapter == "mysql" {
			boilerGeneratorForMysql()
		} else {
			panic(fmt.Errorf("Unknown adapter %s", config.Test.Adapter))
		}
	},
}

func init() {
	rootCmd.AddCommand(sqlboilerCmd)
}
