package spec

import (
	"context"
	"{{.Module}}/config"
	"{{.Module}}/controllers"
	"{{.Module}}/models"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/stretchr/testify/assert"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"
)

func TestMain(m *testing.M) {
	RunTestMain(m)
}

// TestCreate{{.Model}}
// POST /{{.models}}/
func TestCreate{{.Model}}(t *testing.T) {
	{{.model}}JSON := `{
  	{{- range $column := .Table.Columns -}}
    {{if eq $column.Type "String" }}
		"${$column.Name}": "one",
    {{end -}}
    {{if eq $column.Type "Text" }}
		"${$column.Name}": "two",
    {{end -}}
    {{if eq $column.Type "Integer" }}
		"${$column.Name}": 1,
    {{end -}}
    {{if eq $column.Type "Boolean" }}
		"${$column.Name}": true,
    {{end -}}
    {{end -}}
	}`

	t.Run("Success", func(t *testing.T) {
		{{.models}}Count := count{{.models}}()
		c, rec := Request(http.MethodPost, "/{{.models}}", {{.model}}JSON)
		{{.models}}Controller := &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(t,
			config.RouteTo({{.models}}Controller, {{.models}}Controller.Create)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"first_name": "Jon Snow",
			"last_name": "James",
			"email": "jon@labstack.com",
			"id": "<<PRESENCE>>",
			"encrypted_password": "<<PRESENCE>>",
			"created_at": "<<PRESENCE>>",
			"updated_at": "<<PRESENCE>>"
		}`)
		assert.Equal(t, {{.models}}Count+1, count{{.models}}())
	})

	t.Run("Failure", func(t *testing.T) {
		{{.models}}Count := count{{.models}}()
		c, rec := Request(http.MethodPost, "/{{.models}}", `{wrongJSON}`)
		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(
			t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Create)(c),
		)

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "<<PRESENCE>>"
		}`)
		assert.Equal(t, {{.models}}Count, count{{.models}}())
	})
}

// TestGet{{.Model}}
// GET /{{.models}}/:id
func TestGet{{.Model}}(t *testing.T) {
	BeforeEachTest()

	t.Run("{{.Model}} Exist", func(t *testing.T) {
		c, rec := Request(http.MethodGet, "/{{.models}}/1")
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		ja := jsonassert.New(t)
		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(
			t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Get)(c))

		ja.Assertf(rec.Body.String(), `{
			"first_name": "James",
			"last_name": "Huynh",
			"email": "james@rubify.com",
			"id": 1,
			"encrypted_password": "<<PRESENCE>>",
			"created_at": "<<PRESENCE>>",
			"updated_at": "<<PRESENCE>>"
		}`)
	})

	t.Run("{{.Model}} does not exist", func(t *testing.T) {
		c, rec := Request(http.MethodGet, "/{{.models}}/1")
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		ja := jsonassert.New(t)

		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Get)(c))
		ja.Assertf(rec.Body.String(), `{
			"error": "{{.Model}} does not exist"
		}`)
	})
}

// TestUpdate{{.Model}}
// PATCH /{{.models}}/:id
func TestUpdate{{.Model}}(t *testing.T) {
	{{.model}}JSON := `{
		"first_name":"Updated First Name",
		"last_name": "Updated Last Name",
		"email":"james+updated@rubify.com"
	}`

	BeforeEachTest()
	t.Run("{{.Model}} Exist", func(t *testing.T) {
		c, rec := Request(http.MethodPatch, "/{{.models}}/1", {{.model}}JSON)
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Update)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"first_name": "Updated First Name",
			"last_name": "Updated Last Name",
			"email": "james+updated@rubify.com",
			"id": 1,
			"encrypted_password": "<<PRESENCE>>",
			"created_at": "<<PRESENCE>>",
			"updated_at": "<<PRESENCE>>"
		}`)
	})

	t.Run("{{.Model}} does not exist", func(t *testing.T) {
		c, rec := Request(http.MethodPut, "/{{.models}}/0", {{.model}}JSON)
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Update)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "{{.Model}} does not exist"
		}`)
	})
}

// TestDelete{{.Model}}
// DELETE /{{.models}}/:id
func TestDelete{{.Model}}(t *testing.T) {
	BeforeEachTest()

	t.Run("{{.Model}} Exist", func(t *testing.T) {
		{{.models}}Count := count{{.models}}()
		c, _ := Request(http.MethodDelete, "/{{.models}}/1")
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Delete)(c))
		assert.Equal(t, {{.models}}Count-1, count{{.models}}())
	})

	t.Run("{{.Model}} does not exist", func(t *testing.T) {
		{{.models}}Count := count{{.models}}()
		c, rec := Request(http.MethodDelete, "/{{.models}}/0")
		c.SetPath("/{{.models}}/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(
			t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Delete)(c))
		assert.Equal(t, {{.models}}Count, count{{.models}}())
		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "{{.Model}} does not exist"
		}`)
	})
}

// TestIndex{{.models}}
// GET /{{.models}}/
func TestIndex{{.models}}(t *testing.T) {
	t.Run("{{.Model}} Exist", func(t *testing.T) {
		BeforeEachTest()

		c, rec := Request(http.MethodGet, "/{{.models}}")
		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(
			t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Index)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `[{
			"first_name": "James",
			"last_name": "Huynh",
			"email": "james@rubify.com",
			"id": 1,
			"encrypted_password": "<<PRESENCE>>",
			"created_at": "<<PRESENCE>>",
			"updated_at": "<<PRESENCE>>"
		},
		{
			"first_name": "Ivy",
			"last_name": "Pham",
			"email": "ivy@rubify.com",
			"id": 2,
			"encrypted_password": "<<PRESENCE>>",
			"created_at": "<<PRESENCE>>",
			"updated_at": "<<PRESENCE>>"
		}]`)
	})

	t.Run("{{.Model}} does not exist", func(t *testing.T) {
		{{.models}}, _ := models.{{.models}}().AllG(context.Background())
		{{.models}}.DeleteAllG(context.Background())
		c, rec := Request(http.MethodGet, "/{{.models}}")
		var {{.models}}Controller = &controllers.{{.models}}Controller{EchoContext: c}
		assert.NoError(
			t, config.RouteTo({{.models}}Controller, {{.models}}Controller.Index)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `[]`)
	})
}

func count{{.models}}() int64 {
	{{.models}}Count, _ := models.{{.models}}().CountG(context.Background())
	return {{.models}}Count
}
