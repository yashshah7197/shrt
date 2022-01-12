// Package testgroup maintains the group of handlers for the test routes.
package testgroup

import (
	"context"
	"net/http"

	"github.com/yashshah7197/shrt/foundation/web"

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

	h.Logger.Infow("test", "statusCode", http.StatusOK, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	return web.Respond(ctx, w, status, http.StatusOK)
}
