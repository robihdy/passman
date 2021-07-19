package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/robihdy/passman/internal/data"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/logins", app.requirePermission(data.PermissionCodeLogins, app.listLoginsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/logins", app.requirePermission(data.PermissionCodeLogins, app.createLoginHandler))
	router.HandlerFunc(http.MethodGet, "/v1/logins/:id", app.requirePermission(data.PermissionCodeLogins, app.showLoginHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/logins/:id", app.requirePermission(data.PermissionCodeLogins, app.updateLoginHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/logins/:id", app.requirePermission(data.PermissionCodeLogins, app.deleteLoginHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.authenticate(router)))
}
