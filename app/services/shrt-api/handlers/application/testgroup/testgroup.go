// Package testgroup maintains the group of handlers for the test routes.
package testgroup

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Logger *zap.SugaredLogger
}

// Test is a basic handler for development purposes.
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK
	json.NewEncoder(w).Encode(status)

	h.Logger.Infow("test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}
