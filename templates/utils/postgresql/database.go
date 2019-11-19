package utils

import (
	"database/sql"
	"fmt"

	"github.com/integralist/go-findroot/find"
	"github.com/volatiletech/sqlboiler/boil"

	// import driver
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver"
)

// GetDB return the global DB object
func GetDB() *sql.DB {
	return boil.GetDB().(*sql.DB)
}

// OpenDB open database
func OpenDB(env string) func() {
	basePath, _ := find.Repo()
	config, err := LoadDBConfig(basePath.Path + "/config/database.toml")
	if err != nil {
		panic(err)
	}

	var dbConfig postgresql

	if env == "development" {
		dbConfig = config.Development
	} else if env == "production" {
		dbConfig = config.Production
	} else if env == "staging" {
		dbConfig = config.Staging
	} else if env == "test" {
		dbConfig = config.Test
	}

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password,
		dbConfig.Host, dbConfig.Port, dbConfig.Dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	boil.SetDB(db)

	return func() {
		db.Close()
	}
}
