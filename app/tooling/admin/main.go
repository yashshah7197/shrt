package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

func main() {
	// Generate a new private/public key pair.
	err := genKeyPair()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate a new, signed JSON Web Token.
	err = genToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// genKeyPair generates a new x509 private/public keypair for auth tokens.
func genKeyPair() error {
	// Generate a new private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create a file for the private key information in PEM form.
	privateKeyFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("creating private pem file: %w", err)
	}
	defer privateKeyFile.Close()

	// Construct a PEM block for the private key.
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Write the private key to the private key file.
	if err := pem.Encode(privateKeyFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private key file: %w", err)
	}

	// Marshal the public key from the private key to PKIX.
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	// Create a file for the public key information in PEM form.
	publicKeyFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("creating public pem file: %w", err)
	}
	defer publicKeyFile.Close()

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	// Write the public key to the public key file.
	if err := pem.Encode(publicKeyFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to public key file: %w", err)
	}

	return nil
}

// genToken generates a new, signed JSON Web Token.
func genToken() error {
	// Open the private key file for reading.
	privateFile, err := os.Open("private.pem")
	if err != nil {
		return fmt.Errorf("opening private key file: %w", err)
	}

	// Read the contents of the private key file.
	privateBytes, err := io.ReadAll(privateFile)
	if err != nil {
		return fmt.Errorf("reading private key file: %w", err)
	}

	// Decode the contents of the private key file in to a PEM block.
	privatePEM, _ := pem.Decode(privateBytes)

	// Parse the private key PEM block in to a private key.
	privateKey, err := x509.ParsePKCS1PrivateKey(privatePEM.Bytes)
	if err != nil {
		return fmt.Errorf("parsing private key pem block: %w", err)
	}

	// Generate a new JSON Web Token.
	token, err := jwt.NewBuilder().
		Issuer("shrt-api").
		Subject("b0ef2788-614a-47b6-a7ba-c0c7c75f6d7f").
		IssuedAt(time.Now().UTC()).
		Expiration(time.Now().Add(8760*time.Hour).UTC()).
		Claim("roles", []string{"admin"}).
		Build()
	if err != nil {
		return fmt.Errorf("generating jwt: %w", err)
	}

	// Add the key ID header.
	headers := jws.NewHeaders()
	err = headers.Set("kid", "905789cb-61d7-44c8-a7e2-0e51d43d8c85")
	if err != nil {
		return fmt.Errorf("setting jwt headers: %w", err)
	}

	// Sign the token with our private key.
	signedToken, err := jwt.Sign(token, jwa.RS256, privateKey, jwt.WithHeaders(headers))
	if err != nil {
		return fmt.Errorf("signing token with the private key: %w", err)
	}

	// Print the signed token to STDOUT.
	fmt.Println(string(signedToken))

	return nil
}
