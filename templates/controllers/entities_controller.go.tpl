package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"{{.Module}}/models"
	"{{.Module}}/utils"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo"
	"github.com/volatiletech/sqlboiler/boil"
)

// {{.Models}Controller handle actions on {{.models}}
type {{.Models}}Controller struct {
	{{.model}}        *models.{{.Model}}
	EchoContext echo.Context
}

////////////////////////////// GENERAL METHODS ////////////////////////////////

// BeforeAction is used to handle all BeforeAction
// This method must be implemented
func (uc *{{.Models}}Controller) BeforeAction(
	httpMethod string,
	actionName string,
	c echo.Context,
) error {
	if utils.Contains([]string{"Get", "Update", "Delete"}, actionName) {
		return uc.retrieve{{.Model}}()
	}
	return nil
}

// RenderJSON is used to RenderJSON
// This method must be implemented
func (uc *{{.Models}}Controller) RenderJSON(json interface{}) error {
	_, isError := json.(error)
	if isError {
		return utils.HandleErrorJSON(uc.EchoContext, json.(error))
	}
	return uc.EchoContext.JSON(http.StatusOK, json)
}

////////////////////////////// GENERAL METHODS ENDS ///////////////////////////

////////////////////////////////// ACTIONS START //////////////////////////////

// Index handles return list of {{.model}}s
// GET /{{.model}}s
func (uc *{{.Models}}Controller) Index() error {
	{{.model}}s, err := models.{{.Models}}().AllG(context.Background())
	if err != nil {
		return uc.RenderJSON(err)
	}
	if len({{.model}}s) == 0 {
		return uc.RenderJSON([]string{})
	}
	return uc.RenderJSON({{.model}}s)
}

// Create create {{.model}} using echo context
// POST /{{.model}}s
func (uc *{{.Models}}Controller) Create() error {
	uc.{{.model}} = &models.{{.Model}}{}
	if err := uc.EchoContext.Bind(uc.{{.model}}); err != nil {
		return uc.RenderJSON(err)
	}

	if err := uc.{{.model}}.InsertG(
		context.Background(),
		whiteListColumns(),
	); err != nil {
		return uc.RenderJSON(err)
	}

	return uc.RenderJSON(uc.{{.model}})
}

// Get handles return a {{.model}}
// GET /{{.model}}s/:id
func (uc *{{.Models}}Controller) Get() error {
	return uc.RenderJSON(uc.{{.model}})
}

// Update handles update a {{.model}}
// PATCH /{{.model}}s/:id
func (uc *{{.Models}}Controller) Update() error {
	if err := uc.EchoContext.Bind(uc.{{.model}}); err != nil {
		return uc.RenderJSON(err)
	}

	if _, err := uc.{{.model}}.UpdateG(
		context.Background(),
		whiteListColumns(),
	); err != nil {
		return uc.RenderJSON(err)
	}

	return uc.RenderJSON(uc.{{.model}})
}

// Delete handles delete a {{.model}}
// DELETE /{{.model}}s/:id
func (uc *{{.Models}}Controller) Delete() error {
	uc.{{.model}}.DeleteG(context.Background())
	return uc.EchoContext.NoContent(http.StatusNoContent)
}

////////////////////////////////// ACTIONS END// //////////////////////////////

////////////////////////////////// PRIVATE ////////////////////////////////////

func (uc *{{.Models}}Controller) retrieve{{.Model}}() error {
	id, _ := strconv.Atoi(uc.EchoContext.Param("id"))
	{{.model}}, err := models.Find{{.Model}}G(context.Background(), int(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("{{.Model}} does not exist")
		}
		return err
	}
	uc.{{.model}} = user

	return nil
}

//////////////////////////////// PAREMTERS FILTERER ///////////////////////////

func whiteListColumns() boil.Columns {
	// TODO: must add WhiteList only columns that needed to be updated, i.e.
	// return boil.Whitelist("name", "email")

	return boil.Infer()
}
