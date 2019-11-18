package config

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"{{.AppName}}/utils"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/volatiletech/sqlboiler/boil"
)

// SetupEnvironment setups environment based on environment variable
func SetupEnvironment() (environment string) {
	availableEnvs := []string{"development", "test", "staging", "production"}

	if os.Getenv("GOENV") != "" {
		environment = os.Getenv("GOENV")
	} else {
		environment = "development"
	}
	if environment == "development" || environment == "staging" {
		boil.DebugMode = true
	}
	if utils.Contains(availableEnvs, environment) == false {
		fmt.Printf("Unable to start server in %s environment\n", environment)
		os.Exit(3)
	}
	return
}

func getFunctionName(i interface{}) string {
	fnNames := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	splits := strings.Split(fnNames, ".")
	fnName := strings.ReplaceAll(splits[len(splits)-1], "-fm", "")
	return fnName
}

func getControllerName(controller ApplicationController) string {
	controllerName := reflect.TypeOf(controller).String()
	splits := strings.Split(controllerName, ".")
	return splits[len(splits)-1]
}

func logRequest(
	controller ApplicationController,
	action func() error,
	c echo.Context,
) {
	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString(
		color.YellowString(
			"Processing by %s#%s\n",
			getControllerName(controller),
			getFunctionName(action),
		))

	os.Stdout.Write(buf.Bytes())
}

// RouteTo handle route request
func RouteTo(controller ApplicationController,
	action func() error) func(c echo.Context) error {
	return func(c echo.Context) error {
		if os.Getenv("GOENV") != "test" {
			logRequest(controller, action, c)
		}
		err := controller.BeforeAction(
			c.Request().Method,
			getFunctionName(action),
			c,
		)
		if err != nil {
			return utils.HandleErrorJSON(c, err)
		}
		status := c.Response().Status
		if status == 400 {
			return nil
		}

		err = action()
		return err
	}
}
