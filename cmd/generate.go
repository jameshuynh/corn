package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cobra"

	pluralize "github.com/gertd/go-pluralize"

	"github.com/elliotchance/pie/pie"

	"github.com/iancoleman/strcase"
)

func generateMigration(args []string) string {
	migrationFile := args[1]
	usr, _ := user.Current()
	exec.Command("mkdir", "-p", "db/migrations").CombinedOutput()
	output, _ := exec.Command(
		usr.HomeDir+"/go/bin/mig", "create", migrationFile, "-d", "db/migrations",
	).CombinedOutput()
	fmt.Print(string(output))
	exec.Command("chmod", "-R", "0755", "db/migrations").CombinedOutput()

	return string(output)
}

func generateModel(args []string) {
	if len(args) < 2 {
		panic(errors.New("A model name must be supplied"))
	}

	modelName := args[1]
	rest := args[2:]

	pluralize := pluralize.NewClient()

	pluralModelName := pluralize.Plural(strcase.ToSnake(modelName))

	output := generateMigration([]string{"generate", "create_" + pluralModelName})
	filePath := strings.TrimSpace(strings.Split(string(output), " ")[1])
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(pluralModelName).IfNotExists()
	ctb.Define(
		"id",
		"serial",
		"PRIMARY KEY",
	)

	pie.Strings(rest).Each(func(value string) {
		split := strings.Split(value, ":")
		var (
			fieldType, fieldName string
		)
		if len(split) > 1 {
			fieldType = split[1]
			fieldName = split[0]
		} else {
			fieldType = "string"
			fieldName = split[0]
		}

		if fieldType == "string" {
			ctb.Define(fieldName, "VARCHAR(255)")
		} else if fieldType == "text" {
			ctb.Define(fieldName, "TEXT")
		} else if fieldType == "integer" {
			ctb.Define(fieldName, "INTEGER")
		} else if fieldType == "float" {
			ctb.Define(fieldName, "NUMERIC(2, 4)")
		} else if fieldType == "boolean" {
			ctb.Define(fieldName, "BOOLEAN")
		} else {
			panic(fmt.Errorf("Type %s is not yet supported", fieldType))
		}
	})

	ctb.Define(
		"created_at",
		"TIMESTAMP",
		"DEFAULT",
		"CURRENT_TIMESTAMP",
	)
	ctb.Define(
		"updated_at",
		"TIMESTAMP",
		"DEFAULT",
		"CURRENT_TIMESTAMP",
	)

	cmd := fmt.Sprintf("echo \"%s\" | sql-formatter-cli", ctb.String())
	createSQL, err := exec.Command(
		"bash",
		"-c",
		cmd,
	).CombinedOutput()

	data := fmt.Sprintf(`-- +mig Up
%s;

-- +mig Down
DROP TABLE %s;`, createSQL, pluralModelName)

	err = ioutil.WriteFile(filePath, []byte(data), 0755)
	if err != nil {
		panic(err)
	}
}

var generateCmd = &cobra.Command{
	Use:   "g",
	Short: "Generate Different Types",
	Long:  "Generate Different Types",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			panic(errors.New("No supplied arguments"))
		}

		if args[0] == "migration" {
			generateMigration(args)
		} else if args[0] == "model" {
			generateModel(args)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
