// Package handlers contains the full set of handler functions and routes supported by the API.
package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
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
