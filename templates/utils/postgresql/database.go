package utils

import (
	"database/sql"
	"fmt"
)

// OpenDB open database
func OpenDB(env string) *sql.DB {
	config, err := LoadDBConfig("./config/sqlboiler.toml")
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
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	return db
}
