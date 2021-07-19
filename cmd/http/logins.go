package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/robihdy/passman/internal/data"
	"github.com/robihdy/passman/internal/validator"
	"gopkg.in/guregu/null.v4"
)

func (app *application) createLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string      `json:"name"`
		Username string      `json:"username"`
		Password string      `json:"password"`
		Website  null.String `json:"website"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	login := &data.Login{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
		Website:  input.Website,
	}

	v := validator.New()

	if data.ValidateLogin(v, login); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Logins.Insert(login, app.contextGetUser(r).ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/logins/%d", login.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"login": login}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showLoginHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	login, err := app.models.Logins.Get(id, app.contextGetUser(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"login": login}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateLoginHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	login, err := app.models.Logins.Get(id, app.contextGetUser(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Password *string `json:"password"`
		Website  *string `json:"website"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		login.Name = *input.Name
	}
	if input.Username != nil {
		login.Username = *input.Username
	}
	if input.Password != nil {
		login.Password = *input.Password
	}
	if input.Website != nil {
		login.Website = null.StringFrom(*input.Website)
	}

	v := validator.New()
	if data.ValidateLogin(v, login); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Logins.Update(login)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"login": login}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteLoginHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Logins.Delete(id, app.contextGetUser(r).ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "login successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listLoginsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Username string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Username = app.readString(qs, "username", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "username", "website", "-id", "-name", "-username", "-website"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	logins, err := app.models.Logins.GetByUserID(app.contextGetUser(r).ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"logins": logins}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
