package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// DBConfig Type for MySql
type DBConfig struct {
	Development postgresql
	Test        postgresql
	Staging     postgresql
	Production  postgresql
}

// Psql
type postgresql struct {
	Dbname   string
	Host     string
	Port     int64
	User     string
	Password string
	Sslmode  string
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
