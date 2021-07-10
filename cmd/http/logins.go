package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/robihdy/passman/internal/data"
	"github.com/robihdy/passman/internal/encryption"
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
		Password: encryption.Encrypt(input.Password, app.config.encryption.masterKey),
		Website:  input.Website,
	}

	v := validator.New()

	if data.ValidateLogin(v, login); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Logins.Insert(login)
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

	login, err := app.models.Logins.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	login.Password = encryption.Decrypt(login.Password, app.config.encryption.masterKey)

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

	login, err := app.models.Logins.Get(id)
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
		Name     string      `json:"name"`
		Username string      `json:"username"`
		Password string      `json:"password"`
		Website  null.String `json:"website"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	login.Name = input.Name
	login.Username = input.Username
	login.Password = encryption.Encrypt(input.Password, app.config.encryption.masterKey)
	login.Website = input.Website

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

	login.Password = encryption.Decrypt(login.Password, app.config.encryption.masterKey)

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

	err = app.models.Logins.Delete(id)
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
