package middleware

import (
	"context"
	"errors"
	"fmt"
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

// Authorize validates that an authenticated user has at least one role from a specified list.
func Authorize(roles ...string) web.Middleware {
	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Ensure that the claims are present in the context.
			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for that action"),
					http.StatusForbidden,
				)
			}

			// Check that the claims have one of the authorized roles.
			if !claims.Authorized(roles...) {
				return validate.NewRequestError(
					fmt.Errorf("you are not authorized for that action"),
					http.StatusForbidden,
				)
			}

			// Call the next handler.
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
