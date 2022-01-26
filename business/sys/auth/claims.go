package auth

import (
	"context"
	"errors"
	"time"
)

// These are the expected values for Claims.Roles.
const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

// Claims represents the set of authorization claims transmitted via a JWT.
type Claims struct {
	Issuer    string
	Subject   string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Roles     []string
}

// Authorized returns true if the claims has at least one of the provided roles.
func (c Claims) Authorized(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}

	return false
}

// ctxKeyClaims represents the type of value for the context key.
type ctxKeyClaims int

// claimsKey is how Claims values are stored and retrieved.
const claimsKey ctxKeyClaims = 1

// SetClaims stores the Claims in the context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) (Claims, error) {
	claims, ok := ctx.Value(claimsKey).(Claims)
	if !ok {
		return Claims{}, errors.New("claims value missing from context")
	}

	return claims, nil
}
