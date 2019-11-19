package spec

import (
	"database/sql"
	"encoding/json"
	"io"
	"{{.Module}}/utils"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/integralist/go-findroot/find"
	"gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

// Request is used to create a http request
func Request(method string, url string, bodyOptional ...string) (
	echo.Context, *httptest.ResponseRecorder,
) {

	var body io.Reader
	if len(bodyOptional) > 0 {
		body = strings.NewReader(bodyOptional[0])
	}

	e := echo.New()
	req := httptest.NewRequest(
		method,
		url,
		body,
	)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// JSON convert ResponseRecorder to JSON Object
func JSON(rec *httptest.ResponseRecorder) map[string]interface{} {
	responseMap := make(map[string]interface{})
	json.Unmarshal([]byte(rec.Body.String()), &responseMap)
	return responseMap
}

// JSONArray convert ResponsRecorder to JSONArray Object
func JSONArray(rec *httptest.ResponseRecorder) []map[string]interface{} {
	responseMap := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(rec.Body.String()), &responseMap)
	return responseMap
}

func setupFixtures() {
	basePath, _ := find.Repo()
	fixtures, _ = testfixtures.NewFolder(
		GetDB(), &testfixtures.PostgreSQL{}, basePath.Path+"/spec/fixtures",
	)
}

// RunTestMain should be called in TestMain
func RunTestMain(m *testing.M) {
	os.Setenv("GOENV", "test")
	close := utils.OpenDB("test")
	defer close()
	setupFixtures()
	os.Exit(m.Run())
}

// BeforeEachTest could be called before each test
func BeforeEachTest() {
	fixtures.Load()
}

// GetDB is used to get DB
func GetDB() *sql.DB {
	return boil.GetDB().(*sql.DB)
}

// ConvertToInt converts interpreted id in float to int
func ConvertToInt(id interface{}) int {
	return int(id.(float64))
}
