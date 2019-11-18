package main

import (
	"database/sql"
	"{{.AppName}}/utils"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/volatiletech/sqlboiler/boil"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver"
)

var db *sql.DB

func main() {
	e := echo.New()

  env := os.Getenv("GOENV")
	if env == "" {
		env = "development"
	}

	db = utils.OpenDB(env)
	boil.SetDB(db)
	defer db.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete,
		},
	}))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusCreated, "API Server is online")
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
