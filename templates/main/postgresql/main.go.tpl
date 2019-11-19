package main

//go:generate sqlboiler --wipe psql --add-global-variants
import (
	"{{.Module}}/config"
	"{{.Module}}/utils"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver"
)

func setupMiddleware(e *echo.Echo) {
	// Middleware
	e.Use(utils.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete,
		},
	}))

}

func main() {
	e := echo.New()

	environment := config.SetupEnvironment()
	close := utils.OpenDB(environment)
	defer close()

	setupMiddleware(e)
	config.Routes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
