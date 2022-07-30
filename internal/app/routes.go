package app

import (
	"fmt"
	"net/http"
)

const (
	urlVarUserID = "userID"
)

func (app *App) routes() {
	// handles routing app functionality
	appRouter := app.router.PathPrefix("").Subrouter()

	if app.env != "production" {
		fs := http.FileServer(http.Dir("./api/OpenAPI/"))
		appRouter.PathPrefix("/docs").Handler(http.StripPrefix("/docs", fs))
	}

	appRouter.HandleFunc("/auth", app.handleAuthenticate()).Methods(http.MethodPost)
	appRouter.HandleFunc("/ping", app.handlePing()).Methods(http.MethodGet)

	// handles routing domain functionality for api v1
	apiV1Router := app.router.PathPrefix("/api/v1").Subrouter()
	// Users
	apiV1Router.HandleFunc("/users", app.handleCreateUser()).Methods(http.MethodPost)
	// apiV1Router.HandleFunc("/users", app.handleGetUsers()).Methods("GET")
	apiV1Router.HandleFunc(fmt.Sprintf("/users/{%s}", urlVarUserID), app.handleGetUser()).Methods(http.MethodGet)
	apiV1Router.HandleFunc(fmt.Sprintf("/users/{%s}", urlVarUserID), app.handlePatchUser()).Methods(http.MethodPatch)
	// apiV1Router.HandleFunc(fmt.Sprintf("/users/{%s}", app.handleDeleteUser()).Methods("DELETE")
	// apiV1Router.HandleFunc(fmt.Sprintf("/users/{%s}/password", app.handleUpdateUserPassword()).Methods("PUT")

	apiV1Router.Use(NewRequestResponseLogger(app.logger).Middleware)
	if app.env == "production" || app.env == "staging" {
		//do summat
	}
}
