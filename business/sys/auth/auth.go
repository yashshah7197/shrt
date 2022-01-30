// Package auth provides authentication and authorization support.
package auth

import (
	"fmt"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/yashshah7197/shrt/foundation/keystore"
)

// Auth is used to authenticate clients. It can generate a token for a set of user claims and
// recreate the claims by parsing a token.
type Auth struct {
	activeKeyID        string
	keystore           keystore.KeyStore
	signatureAlgorithm jwa.SignatureAlgorithm
}

// New creates a new Auth to support authentication and authorization.
func New(activeKeyID string, keystore keystore.KeyStore) (*Auth, error) {
	// The activeKeyID represents the private key used to sign new tokens.
	_, err := keystore.PrivateKey(activeKeyID)
	if err != nil {
		return nil, fmt.Errorf("looking up private key: %w", err)
	}

	a := Auth{
		activeKeyID:        activeKeyID,
		keystore:           keystore,
		signatureAlgorithm: jwa.RS256,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	// Generate a new JSON Web Token.
	token, err := jwt.NewBuilder().
		Issuer(claims.Issuer).
		Subject(claims.Subject).
		IssuedAt(claims.IssuedAt).
		Expiration(claims.ExpiresAt).
		Claim("roles", claims.Roles).
		Build()
	if err != nil {
		return "", fmt.Errorf("generating token: %w", err)
	}

	// Fetch the private key associated with the active key id from the keystore.
	privateKey, err := a.keystore.PrivateKey(a.activeKeyID)
	if err != nil {
		return "", fmt.Errorf("fetching private key: %w", err)
	}

	// Sign the token with the private key associated with the active key id.
	signedToken, err := jwt.Sign(token, a.signatureAlgorithm, privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token with private key: %w", err)
	}

	return string(signedToken), nil
}

// ValidateToken parses a JSON Web Token, verifies it and then validates it.
func (a *Auth) ValidateToken(tokenString string) (Claims, error) {
	// Fetch the public key associated with the active key id from the keystore.
	publicKey, err := a.keystore.PublicKey(a.activeKeyID)
	if err != nil {
		return Claims{}, fmt.Errorf("fetching public key: %w", err)
	}

	// Verify and validate the token.
	token, err := jwt.ParseString(
		tokenString,
		jwt.WithVerify(a.signatureAlgorithm, publicKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("validating token: %w", err)
	}

	// Parse the roles from the token claims
	roles, ok := token.Get("roles")
	if !ok {
		return Claims{}, fmt.Errorf("parsing roles from token claims: %w", err)
	}

	// Recreate the claims from the token.
	claims := Claims{
		Issuer:    token.Issuer(),
		Subject:   token.Subject(),
		IssuedAt:  token.IssuedAt(),
		ExpiresAt: token.Expiration(),
		Roles:     roles.([]string),
	}

	return claims, nil
}
