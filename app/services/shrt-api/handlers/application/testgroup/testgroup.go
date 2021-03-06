// Package testgroup maintains the group of handlers for the test routes.
package testgroup

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/yashshah7197/shrt/business/sys/validate"
	"github.com/yashshah7197/shrt/foundation/web"

	"go.uber.org/zap"
)

// Handlers manages the set of check endpoints.
type Handlers struct {
	Logger *zap.SugaredLogger
}

// Test is a basic handler for development purposes.
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return validate.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
