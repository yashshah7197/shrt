// Package web contains a small web framework extension.
package web

import (
	"os"
	"syscall"

	"github.com/go-chi/chi/v5"
)

// App is the entrypoint into our application and what configures our context object for each of
// our HTTP handlers.
type App struct {
	*chi.Mux
	shutdown chan os.Signal
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	return &App{
		Mux:      chi.NewMux(),
		shutdown: shutdown,
	}
}

// SignalShutdown is used to gracefully shut down the application when an integrity issue is
// identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}
