package app

import (
	"fmt"
	"net/http"

	"github.com/OJOMB/graffiti-berlin-svc/internal/app/middleware"
)

const (
	urlVarUserID = "userID"
)

func (app *App) routes() {
	if app.docsEnabled {
		fs := http.FileServer(http.Dir("./api/OpenAPI/"))
		app.router.PathPrefix("/docs").Handler(http.StripPrefix("/docs", fs))
	}

	// Users
	app.router.HandleFunc("/api/v1/users", app.handleCreateUser()).Methods(http.MethodPost)
	// app.router.HandleFunc("/users", app.handleGetUsers()).Methods("GET")
	app.router.HandleFunc(fmt.Sprintf("/api/v1/users/{%s}", urlVarUserID), app.handleGetUser()).Methods(http.MethodGet)
	app.router.HandleFunc("/users/{id}", app.handlePatchUser()).Methods(http.MethodPatch)
	// app.router.HandleFunc("/users/{id}", app.handleDeleteUser()).Methods("DELETE")
	// app.router.HandleFunc("/users/{id}/password", app.handleUpdateUserPassword()).Methods("PUT")

	app.router.Use(middleware.NewRequestResponseLogger(app.logger).Middleware)
}
