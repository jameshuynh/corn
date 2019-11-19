package utils

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// Contains check a string is inside an array of string
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// HandleErrorJSON will convert JSON for error message
func HandleErrorJSON(c echo.Context, err error) error {
	errorMsg := make(map[string]string)
	errorString := err.Error()
	if strings.Contains(errorString, "code=400") {
		errorMsg["error"] = strings.Split(errorString, "message=")[1]
	} else {
		errorMsg["error"] = errorString
	}
	return c.JSON(http.StatusBadRequest, errorMsg)
}
