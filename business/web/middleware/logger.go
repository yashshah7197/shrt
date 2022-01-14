package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/yashshah7197/shrt/foundation/web"

	"go.uber.org/zap"
)

// Logger writes some information about the request to the logs.
func Logger(logger *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// If the context is missing our web values, return an error so that it can be handled
			// further up the chain.
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			logger.Infow(
				"request started",
				"traceid", v.TraceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remoteaddr", r.RemoteAddr,
			)

			// Call the next handler.
			err = handler(ctx, w, r)

			logger.Infow(
				"request completed",
				"traceid", v.TraceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remoteaddr", r.RemoteAddr,
				"statuscode", v.StatusCode,
				"since",
				time.Since(v.Now),
			)

			// Return the error so that it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
