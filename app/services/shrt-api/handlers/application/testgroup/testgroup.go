// Package testgroup maintains the group of handlers for the test routes.
package testgroup

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Logger *zap.SugaredLogger
}

// Test is a basic handler for development purposes.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	statusCode := http.StatusOK
	h.Logger.Infow("test", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return json.NewEncoder(w).Encode(status)
}
