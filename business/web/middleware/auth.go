package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/yashshah7197/shrt/business/sys/auth"
	"github.com/yashshah7197/shrt/business/sys/validate"
	"github.com/yashshah7197/shrt/foundation/web"
)

// Authenticate validates a JSON Web Token from the 'Authorization' header.
func Authenticate(a *auth.Auth) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Expecting: bearer <token>
			authString := r.Header.Get("Authorization")

			// Parse the 'Authorization' header.
			parts := strings.Split(authString, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: bearer <token>")
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// Validate that the token was signed by us.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			// Add the claims to the context so that they can be retrieved later.
			ctx = auth.SetClaims(ctx, claims)

			// Call the next handler.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
