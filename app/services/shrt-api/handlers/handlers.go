// Package handlers contains the full set of handler functions and routes supported by the API.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/yashshah7197/shrt/app/services/shrt-api/handlers/application/testgroup"
	"github.com/yashshah7197/shrt/app/services/shrt-api/handlers/debug/checkgroup"
	"github.com/yashshah7197/shrt/business/web/middleware"
	"github.com/yashshah7197/shrt/foundation/web"

	"go.uber.org/zap"
)

// DebugStandardLibraryMux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServeMux. Using the DefaultServeMux could be a security risk
// since a dependency could inject a handler into our service without us knowing it.
func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

// DebugMux registers all the debug standard library routes and then custom debug application
// routes for the service. This bypassing the use of the DefaultServeMux. Using the DefaultServeMux
// could be a security risk since a dependency could inject a handler into our service without
// us knowing it.
func DebugMux(build string, logger *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// Register the debug check endpoints.
	cgh := checkgroup.Handlers{
		Build:  build,
		Logger: logger,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	return mux
}

// APIMuxConfig contains all the mandatory systems required by the handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Logger   *zap.SugaredLogger
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	// Construct the web.App which holds all the routes.
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Logger),
		middleware.Errors(cfg.Logger),
	)

	// Bind the different routes for the API.
	bindRoutes(app, cfg)

	return app
}

// bindRoutes binds all the API routes to their handlers.
func bindRoutes(app *web.App, cfg APIMuxConfig) {
	tgh := testgroup.Handlers{
		Logger: cfg.Logger,
	}

	app.Handle(http.MethodGet, "/test", tgh.Test)
}
