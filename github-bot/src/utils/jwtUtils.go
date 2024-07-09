package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// DecodePrivateKey decodes a base64 encoded RSA private key and returns the parsed rsa.PrivateKey.
func DecodePrivateKey(privateKeyBase64 string) (*rsa.PrivateKey, error) {
	// Decode the base64 encoded private key
	privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding private key: %v", err)
	}

	// Decode the PEM block containing the private key
	block, _ := pem.Decode(privateKeyDecoded)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Parse the RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	return privateKey, nil
}

// GenerateJWT generates a JSON Web Token (JWT) signed with the given RSA private key for a GitHub App.
func GenerateJWT(appID int64, key *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute * 10).Unix(),
		Issuer:    fmt.Sprintf("%d", appID),
	}

	// Create a new JWT token with the specified claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the RSA private key
	return token.SignedString(key)
}
