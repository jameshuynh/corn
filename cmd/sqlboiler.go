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
	replacers["GetString(\"psql."] = "GetString(\"test."
	helpers.SearchAndReplaceFiles(".", replacers)
}

func boilerGeneratorForMysql() {
	// TODO
}

func boilerGeneratorForPsql() {
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
	if err != nil {
		panic(err)
	}
	exec.Command(
		"chmod", "-R", "0755", "config/sqlboiler_test.toml",
	).CombinedOutput()
	output, err := exec.Command(
		"sqlboiler", "-c", "config/sqlboiler_test.toml", "--wipe", "psql",
	).CombinedOutput()
	fmt.Println(string(output), err)
	exec.Command("rm", "-rf", "config/sqlboiler_test.toml").CombinedOutput()

	searchAndReplaceSQLBoilerConfig()
}

var sqlboilerCmd = &cobra.Command{
	Use:   "sqlboiler",
	Short: "Generate SQLBoiler",
	Long:  "Generate SQLBoiler",
	Run: func(cmd *cobra.Command, args []string) {
		config, _ := helpers.LoadDBConfig("./config/sqlboiler.toml")
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
