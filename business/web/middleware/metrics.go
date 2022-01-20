package middleware

import (
	"context"
	"net/http"

	"github.com/yashshah7197/shrt/business/sys/metrics"
	"github.com/yashshah7197/shrt/foundation/web"
)

// Metrics updates the application metrics.
func Metrics() web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Add the metrics into the context for metrics gathering.
			ctx = metrics.Set(ctx)

			// Call the next handler.
			err := handler(ctx, w, r)

			// Increment the requests and the goroutines metrics.
			metrics.AddGoroutine(ctx)
			metrics.AddRequest(ctx)

			// If there is an error flowing through the request then increment the errors metric.
			if err != nil {
				metrics.AddError(ctx)
			}

			// Return the error so that it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
