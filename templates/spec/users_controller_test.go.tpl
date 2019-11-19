package spec

import (
	"context"
	"myapp/config"
	"myapp/controllers"
	"myapp/models"
	"net/http"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/stretchr/testify/assert"
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"
)

func TestMain(m *testing.M) {
	RunTestMain(m)
}

// TestCreateUser
// POST /users/
func TestCreateUser(t *testing.T) {
	userJSON := `{
		"first_name":"Jon Snow",
		"last_name": "James",
		"email":"jon@labstack.com"
	}`

	t.Run("Success", func(t *testing.T) {
		usersCount := countUsers()
		c, rec := Request(http.MethodPost, "/users", userJSON)
		usersController := &controllers.UsersController{EchoContext: c}
		assert.NoError(t,
			config.RouteTo(usersController, usersController.Create)(c))

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
		assert.Equal(t, usersCount+1, countUsers())
	})

	t.Run("Failure", func(t *testing.T) {
		usersCount := countUsers()
		c, rec := Request(http.MethodPost, "/users", `{wrongJSON}`)
		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(
			t, config.RouteTo(usersController, usersController.Create)(c),
		)

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "<<PRESENCE>>"
		}`)
		assert.Equal(t, usersCount, countUsers())
	})
}

// TestGetUser
// GET /users/:id
func TestGetUser(t *testing.T) {
	BeforeEachTest()

	t.Run("User Exist", func(t *testing.T) {
		c, rec := Request(http.MethodGet, "/users/1")
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		ja := jsonassert.New(t)
		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(
			t, config.RouteTo(usersController, usersController.Get)(c))

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

	t.Run("User does not exist", func(t *testing.T) {
		c, rec := Request(http.MethodGet, "/users/1")
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		ja := jsonassert.New(t)

		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(t, config.RouteTo(usersController, usersController.Get)(c))
		ja.Assertf(rec.Body.String(), `{
			"error": "User does not exist"
		}`)
	})
}

// TestUpdateUser
// PATCH /users/:id
func TestUpdateUser(t *testing.T) {
	userJSON := `{
		"first_name":"Updated First Name",
		"last_name": "Updated Last Name",
		"email":"james+updated@rubify.com"
	}`

	BeforeEachTest()
	t.Run("User Exist", func(t *testing.T) {
		c, rec := Request(http.MethodPatch, "/users/1", userJSON)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(t, config.RouteTo(usersController, usersController.Update)(c))

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

	t.Run("User does not exist", func(t *testing.T) {
		c, rec := Request(http.MethodPut, "/users/0", userJSON)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(t, config.RouteTo(usersController, usersController.Update)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "User does not exist"
		}`)
	})
}

// TestDeleteUser
// DELETE /users/:id
func TestDeleteUser(t *testing.T) {
	BeforeEachTest()

	t.Run("User Exist", func(t *testing.T) {
		usersCount := countUsers()
		c, _ := Request(http.MethodDelete, "/users/1")
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(t, config.RouteTo(usersController, usersController.Delete)(c))
		assert.Equal(t, usersCount-1, countUsers())
	})

	t.Run("User does not exist", func(t *testing.T) {
		usersCount := countUsers()
		c, rec := Request(http.MethodDelete, "/users/0")
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("0")

		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(
			t, config.RouteTo(usersController, usersController.Delete)(c))
		assert.Equal(t, usersCount, countUsers())
		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `{
			"error": "User does not exist"
		}`)
	})
}

// TestIndexUsers
// GET /users/
func TestIndexUsers(t *testing.T) {
	t.Run("User Exist", func(t *testing.T) {
		BeforeEachTest()

		c, rec := Request(http.MethodGet, "/users")
		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(
			t, config.RouteTo(usersController, usersController.Index)(c))

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

	t.Run("User does not exist", func(t *testing.T) {
		users, _ := models.Users().AllG(context.Background())
		users.DeleteAllG(context.Background())
		c, rec := Request(http.MethodGet, "/users")
		var usersController = &controllers.UsersController{EchoContext: c}
		assert.NoError(
			t, config.RouteTo(usersController, usersController.Index)(c))

		ja := jsonassert.New(t)
		ja.Assertf(rec.Body.String(), `[]`)
	})
}

func countUsers() int64 {
	usersCount, _ := models.Users().CountG(context.Background())
	return usersCount
}
