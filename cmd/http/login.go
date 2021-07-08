package main

import (
	"fmt"
	"net/http"
	"time"

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

	login := data.Login{
		ID:        id,
		Name:      "Reddit",
		Username:  "rabbithole",
		Password:  "asiap1234",
		Website:   null.NewString("https://reddit.com", true),
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"login": login}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
