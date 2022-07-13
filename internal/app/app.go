package app

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	componentApp = "App"
	appHandler   = "handler"
)

type App struct {
	router      *mux.Router
	logger      *logrus.Entry
	addr        net.Addr
	version     string
	docsEnabled bool

	service Service
}

func New(r *mux.Router, logger *logrus.Logger, addr net.Addr, version string, docsEnabled bool, svc Service) *App {
	return &App{
		router:      r,
		logger:      logger.WithField("component", componentApp),
		addr:        addr,
		version:     version,
		docsEnabled: docsEnabled,
		service:     svc,
	}
}

func (app *App) Run() {
	app.routes()
	app.logger.Fatal(http.ListenAndServe(app.addr.String(), app.router))
}
