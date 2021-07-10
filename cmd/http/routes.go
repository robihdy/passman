package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/logins", app.createLoginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/logins/:id", app.showLoginHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/logins/:id", app.updateLoginHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/logins/:id", app.deleteLoginHandler)

	return router
}
