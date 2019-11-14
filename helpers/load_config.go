package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// DBConfig Type for MySql
type DBConfig struct {
	Development Postgresql
	Test        Postgresql
	Staging     Postgresql
	Production  Postgresql
}

// Postgresql is postgresql struct to expose data of db connection
type Postgresql struct {
	Dbname   string
	Host     string
	Port     int64
	Adapter  string
	User     string
	Password string
	Sslmode  string
}

// RetrieveDBConfig helps to get the DBConfig for an environment
func RetrieveDBConfig(environment string) Postgresql {
	config, err := LoadDBConfig("./config/database.toml")
	if err != nil {
		panic(err)
	}

	if environment == "development" {
		return config.Development
	} else if environment == "test" {
		return config.Test
	} else if environment == "production" {
		return config.Production
	}

	panic(fmt.Errorf("Environment %s is not supported", environment))
}

// GenerateDBConfigString generates the config for migration
func GenerateDBConfigString(environment string) (string, string) {
	config, err := LoadDBConfig("./config/database.toml")
	if err != nil {
		panic(err)
	}

	var dbConfig Postgresql

	if environment == "development" {
		dbConfig = config.Development
	} else if environment == "test" {
		dbConfig = config.Test
	}
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
	), dbConfig.Adapter
}

// LoadDBConfig is used to load database config
func LoadDBConfig(configFile string) (*DBConfig, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist")
	} else if err != nil {
		fmt.Println("hello")
		return nil, err
	}

	var conf DBConfig
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {

		fmt.Println(err)
		return nil, err
	}

	return &conf, nil
}
