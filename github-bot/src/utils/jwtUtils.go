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

func DecodePrivateKey(privateKeyBase64 string) (*rsa.PrivateKey, error) {
    privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKeyBase64)
    if err != nil {
        return nil, fmt.Errorf("error decoding private key: %v", err)
    }

    block, _ := pem.Decode(privateKeyDecoded)
    if block == nil || block.Type != "RSA PRIVATE KEY" {
        return nil, fmt.Errorf("failed to decode PEM block containing private key")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("error parsing private key: %v", err)
    }

    return privateKey, nil
}

func GenerateJWT(appID int64, key *rsa.PrivateKey) (string, error) {
    now := time.Now()
    claims := jwt.StandardClaims{
        IssuedAt:  now.Unix(),
        ExpiresAt: now.Add(time.Minute * 10).Unix(),
        Issuer:    fmt.Sprintf("%d", appID),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(key)
}
