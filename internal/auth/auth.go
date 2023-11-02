package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const Key ctxKey = 1

type Auth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewAuth(privateKey *rsa.PrivateKey, publickey *rsa.PublicKey) (*Auth, error) {
	if privateKey == nil || publickey == nil {
		return nil, errors.New("private and public key cannot be nil")

	}
	return &Auth{
		privateKey: privateKey,
		publicKey:  publickey,
	}, nil

}

// GenerateToken is a method for Auth struct. It generates a new JWT token using the provided claims and
// signs it using the privateKey of the Auth struct it's called upon. If there is an error during signing,
// it returns an error.
func (a *Auth) GenerateToken(claims jwt.RegisteredClaims) (string, error) {
	//NewWithClaims creates a new Token with the specified signing method and claims.
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Signing our token with our private key.
	tokenStr, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token %w", err)
	}

	return tokenStr, nil
}

// ValidateToken is a method for Auth struct. It verifies the provided JWT token using the publicKey of the Auth struct
// it's called upon and returns the parsed claims if the JWT token is valid. If the JWT token is invalid or
// there is an error during parsing, it returns an error.
func (a *Auth) ValidateToken(token string) (jwt.RegisteredClaims, error) {
	var c jwt.RegisteredClaims
	// Parse the token with the registered claims.
	tkn, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("parsing token %w", err)
	}
	// Check if the parsed token is valid.
	if !tkn.Valid {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return c, nil
}
