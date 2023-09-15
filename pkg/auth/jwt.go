package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	*jwt.RegisteredClaims
}

func CreateJWT(pemKey []byte) (string, error) {
	expiration := time.Now().Add(time.Second * 180)

	claims := MyCustomClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "77423",
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(pemKey)

	return token.SignedString(privateKey)
}
