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
			traceID := "00000000-0000-0000-0000-000000000000"
			statusCode := http.StatusOK
			now := time.Now()

			logger.Infow(
				"request started",
				"traceid", traceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remoteaddr", r.RemoteAddr,
			)

			err := handler(ctx, w, r)

			logger.Infow(
				"request completed",
				"traceid", traceID,
				"method", r.Method,
				"path", r.URL.Path,
				"remoteaddr", r.RemoteAddr,
				"statuscode", statusCode,
				"since",
				time.Since(now),
			)

			return err
		}

		return h
	}

	return m
}
