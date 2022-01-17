package middleware

import (
	"context"
	"net/http"

	"github.com/yashshah7197/shrt/business/sys/validate"
	"github.com/yashshah7197/shrt/foundation/web"

	"go.uber.org/zap"
)

// Errors handles the errors coming out of the call chain. It detects normal application errors
// which are used to respond to the client in a uniform way. Unexpected errors (status code >= 500)
// are logged.
func Errors(logger *zap.SugaredLogger) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// If the context is missing our web values, return an error so that it can be handled
			// further up the chain.
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			// Run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {
				// Log the error.
				logger.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				// Build out the error response.
				var er validate.ErrorResponse
				var statusCode int
				switch errorType := validate.Cause(err).(type) {
				case validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: errorType.Error(),
					}
					statusCode = http.StatusBadRequest

				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: errorType.Error(),
					}
					statusCode = errorType.StatusCode

				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					statusCode = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, statusCode); err != nil {
					return err
				}

				// If we receive a shutdown error, return it to the base handler, so it can shut
				// down the service.
				if ok := web.IsShutdownError(err); ok {
					return err
				}
			}

			// The error has been handled, so stop propagating it.
			return nil
		}

		return h
	}

	return m
}
