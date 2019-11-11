package utils

import (
	"database/sql"
	"fmt"
)

// OpenDB open database
func OpenDB() *sql.DB {
	config, err := LoadConfig("./config/sqlboiler.toml")
	if err != nil {
		panic(err)
	}

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Psql.User, config.Psql.Password,
		config.Psql.Host, config.Psql.Port, config.Psql.Dbname)
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
