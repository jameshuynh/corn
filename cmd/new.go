package cmd

/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jameshuynh/corn/helpers"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// corn new path --database=mysql
// corn new path --database=postgresql
// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate an Echo project, with SQLBoiler",
	Long:  `Generate an Echo project, with SQLBoiler`,
	Run: func(cmd *cobra.Command, args []string) {
		database, _ := cmd.Flags().GetString("database")
		if len(args) < 1 {
			helpers.ExitOnError(fmt.Errorf("App Path must be supplied"))
		}

		if database != "mysql" && database != "postgresql" {
			helpers.ExitOnError(
				fmt.Errorf(
					"Database \"%s\" is not supported. Only mysql & postgresql are supported",
					database))
		}

		path := args[0]
		generateProjectFolder(path, database)
	},
}

func copyFiles(database string, appPath string, currDir string) {
	exec.
		Command("mkdir", "-p", "utils").
		CombinedOutput()
	exec.
		Command("mkdir", "-p", "config").
		CombinedOutput()
	exec.
		Command("mkdir", "-p", "db").
		CombinedOutput()

	exec.Command(
		"cp", "-rf",
		fmt.Sprintf("%s/templates/utils/%s/", currDir, database),
		"utils/",
	).CombinedOutput()

	exec.Command(
		"cp", "-rf",
		fmt.Sprintf("%s/templates/main/%s/main.go.tmpl", currDir, database),
		"./main.go",
	).CombinedOutput()

	exec.Command(
		"cp", "-rf",
		fmt.Sprintf("%s/templates/config/%s/", currDir, database),
		"config",
	).CombinedOutput()

	exec.Command(
		"cp", "-rf",
		fmt.Sprintf("%s/templates/db/%s/", currDir, database),
		"db",
	).CombinedOutput()
}

func searchAndReplaceProjectName(projectName string) {
	replacers := map[string]string{
		"{{APP_NAME}}": projectName,
	}
	helpers.SearchAndReplaceFiles(".", replacers)
}

func createDatabase(dbName string, databaseType string, env string) {
	dbName = fmt.Sprintf("%s-%s", dbName, env)
	if databaseType == "postgresql" {
		exec.
			Command("dropdb", dbName).
			CombinedOutput()
		exec.
			Command("createdb", dbName, "--encoding=utf-8").
			CombinedOutput()
		exec.
			Command("psql", "-U", "postgres", "-d", dbName, "-f", "db/database.sql").
			CombinedOutput()
	} else if databaseType == "mysql" {
		// *TODO: do for mysql
		exec.
			Command(fmt.Sprintf("echo \"drop database `%s`\" | mysql -u root -p", dbName)).
			CombinedOutput()
		exec.
			Command(fmt.Sprintf("echo \"create database `%s`\" | mysql -u root -p", dbName)).
			CombinedOutput()
	}

	fmt.Printf("created database %s\n", dbName)
}

func generateProjectFolder(appPath string, database string) {
	pathChunks := strings.Split(appPath, "/")

	projectName := pathChunks[len(pathChunks)-1]

	exec.Command("rm", "-rf", appPath).CombinedOutput()
	err := os.MkdirAll(appPath, 0755)
	helpers.ExitOnError(err)

	currDir, err := helpers.GetLatestBaseFolder()
	helpers.ExitOnError(err)

	err = os.Chdir(appPath)
	helpers.ExitOnError(err)

	c := color.New(color.FgGreen)
	c.Printf("Setup module %v\n", projectName)
	exec.Command("go", "mod", "init", projectName).CombinedOutput()

	fmt.Println()

	c = color.New(color.FgBlue)
	c.Println("Install dependencies:")

	fmt.Println("Getting github.com/labstack/echo...")
	exec.
		Command("go", "get", "-u", "github.com/labstack/echo/...").
		CombinedOutput()

	fmt.Println("Getting github.com/volatiletech/sqlboiler...")
	exec.
		Command("go", "get", "-u", "github.com/volatiletech/sqlboiler/...").
		CombinedOutput()

	if database == "mysql" {
		fmt.Println(
			"Getting github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql...")
		exec.Command("go", "get",
			"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql").
			CombinedOutput()
	} else if database == "postgresql" {
		fmt.Println(
			"Getting github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql...")
		exec.Command("go", "get",
			"github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql").
			CombinedOutput()
	}

	fmt.Println("Getting github.com/BurntSushi/toml...")
	exec.
		Command("go", "get", "github.com/BurntSushi/toml").
		CombinedOutput()

	fmt.Println("Getting github.com/elliotchance/pie...")
	exec.
		Command("go", "get", "-u", "github.com/elliotchance/pie").
		CombinedOutput()

	fmt.Println("Getting github.com/volatiletech/mig...")
	exec.
		Command("go", "get", "-u", "github.com/volatiletech/mig/...").
		CombinedOutput()

	fmt.Println("Getting sql-formatter-cli")
	exec.
		Command("npm", "i", "-g", "sql-formatter-cli").
		CombinedOutput()

	copyFiles(database, appPath, currDir)

	exec.Command("chmod", "-R", "0755", ".").CombinedOutput()
	c = color.New(color.FgGreen)

	searchAndReplaceProjectName(projectName)
	fmt.Println()
	createDatabase(projectName, database, "development")
	createDatabase(projectName, database, "test")

	fmt.Println("Generate SQL Boiler models...")
	exec.
		Command("corn", "sqlboiler").
		CombinedOutput()

	exec.Command("git", "init").CombinedOutput()
	c.Println("\nCompleted!")
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.PersistentFlags().String(
		"database",
		"postgresql",
		"database type, available options: mysql, postgresql",
	)
}
