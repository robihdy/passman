package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/logins", app.requirePermission("logins", app.listLoginsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/logins", app.requirePermission("logins", app.createLoginHandler))
	router.HandlerFunc(http.MethodGet, "/v1/logins/:id", app.requirePermission("logins", app.showLoginHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/logins/:id", app.requirePermission("logins", app.updateLoginHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/logins/:id", app.requirePermission("logins", app.deleteLoginHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.authenticate(app.authenticate(router)))
}
