// Package web contains a small web framework extension.
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi/v5"
)

// A Handler is a type that handles an HTTP request within our own small custom web framework
// extension.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context object for each of
// our HTTP handlers.
type App struct {
	*chi.Mux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handles a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		Mux:      chi.NewMux(),
		shutdown: shutdown,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shut down the application when an integrity issue is
// identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair to the application server
// mux.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// First wrap handler specific middleware around this handler
	handler = wrapMiddleware(mw, handler)

	// Then add the application's general middleware to the handler chain
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Call the wrapped handler functions
		if err := handler(ctx, w, r); err != nil {
			// Error handling
			return
		}
	}

	a.MethodFunc(method, path, h)
}
