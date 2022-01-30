package keystore

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
)

// KeyStore represents an in-memory keystore for authentication and authorization.
type KeyStore struct {
	store jwk.Set
}

// New constructs a new, empty KeyStore.
func New() *KeyStore {
	return &KeyStore{
		store: jwk.NewSet(),
	}
}

// NewFS constructs a new KeyStore based on a set of PEM files rooted inside a directory. The name
// of each PEM file will be used as the key id for that particular key.
func NewFS(fsys fs.FS) (*KeyStore, error) {
	ks := KeyStore{
		store: jwk.NewSet(),
	}

	// This is the function that will be used for walking the directory.
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		// Check if the current directory entry is a directory.
		if dirEntry.IsDir() {
			return nil
		}

		// Check if the current file name extension is .pem.
		if path.Ext(fileName) != ".pem" {
			return nil
		}

		// Open the private key file for reading.
		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening private key file :%w", err)
		}
		defer file.Close()

		// Read the contents of the private key file.
		privateFileBytes, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("reading private key file: %w", err)
		}

		// Decode the contents of the private key file in to a PEM block.
		privatePEM, _ := pem.Decode(privateFileBytes)

		// Parse the private key PEM block in to a private key.
		privateKey, err := x509.ParsePKCS1PrivateKey(privatePEM.Bytes)
		if err != nil {
			return fmt.Errorf("parsing private key pem block: %w", err)
		}

		// Add the private key to the keystore.
		if err := ks.Add(privateKey, strings.TrimSuffix(dirEntry.Name(), ".pem")); err != nil {
			return fmt.Errorf("adding key to the keystore: %w", err)
		}

		return nil
	}

	// Walk the current directory and add all keys to the keystore.
	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return nil, fmt.Errorf("walking keys directory: %w", err)
	}

	return &ks, nil
}

// Add adds a private key with its associated key id to the keystore.
func (ks *KeyStore) Add(privateKey *rsa.PrivateKey, keyID string) error {
	// Create a new JWK private key from the given private key.
	jwkPrivateKey, err := jwk.New(privateKey)
	if err != nil {
		return fmt.Errorf("creating jwk private key: %w", err)
	}

	// Set the "kid" header.
	err = jwkPrivateKey.Set(jwk.KeyIDKey, keyID)
	if err != nil {
		return fmt.Errorf("setting kid header: %w", err)
	}

	// Add it to our JWK key set.
	ks.store.Add(jwkPrivateKey)

	return nil
}

// Remove removes a private key associated with a given key id from the keystore.
func (ks *KeyStore) Remove(keyID string) error {
	// Check if a key with the given id exists in our key set.
	privateKey, ok := ks.store.LookupKeyID(keyID)
	if !ok {
		return errors.New("no key was found with the given key id")
	}

	// Remove the key from the key set.
	ks.store.Remove(privateKey)

	return nil
}

// PrivateKey looks up the keystore for a given key id and returns the corresponding private key.
func (ks *KeyStore) PrivateKey(keyID string) (jwk.Key, error) {
	// Check if a key with the given id exists in our key set.
	key, ok := ks.store.LookupKeyID(keyID)
	if !ok {
		return nil, errors.New("no key was found with the given key id")
	}

	return key, nil
}

// PublicKey looks up the keystore for a given key id and returns the corresponding public key.
func (ks *KeyStore) PublicKey(keyID string) (jwk.Key, error) {
	// Check if a key with the given id exists in our key set.
	key, ok := ks.store.LookupKeyID(keyID)
	if !ok {
		return nil, errors.New("no key was found with the given key id")
	}

	publicKey, err := key.PublicKey()
	if err != nil {
		return nil, errors.New("could not get public key from the private key")
	}

	// Return the public key from the private key.
	return publicKey, nil
}
