package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrEmptySecret = errors.New("jwt secret is empty")

func Generate(userID, secret string, expiresIn time.Duration) (string, error) {
	if secret == "" {
		return "", ErrEmptySecret
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Parse(tokenString, secret string) (string, error) {
	if secret == "" {
		return "", ErrEmptySecret
	}

	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(*jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})) // กันการปลอม token ด้วย algorithm อื่น
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}
