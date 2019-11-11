package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config Type for MySql
type Config struct {
	Psql postgresql
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

// LoadConfig is used to load database config
func LoadConfig(configFile string) (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil, errors.New("Config file does not exist")
	} else if err != nil {
		fmt.Println("hello")
		return nil, err
	}

	var conf Config
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {

		fmt.Println(err)
		return nil, err
	}

	return &conf, nil
}
