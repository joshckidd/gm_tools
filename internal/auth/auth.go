// internal package to handle authentication for the api
// this is very rudimentary at the moment

package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

// hash a password string
func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

// check a password string against an existing hash
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

// create a JWT to be used as a bearer token
func MakeJWT(userName string, tokenSecret string, expiresIn time.Duration) (string, error) {
	res := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "dm_tools",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userName,
	})
	return res.SignedString([]byte(tokenSecret))
}

// validate a provided bearer token
func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}
	sub, err := tok.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	exp, err := tok.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}
	if exp.Time.Before(time.Now()) {
		return "", errors.New("Token expired")
	}

	return sub, nil
}

// extract the bearer token from http headers
func GetBearerToken(headers http.Header) (string, error) {
	res := strings.Split(headers.Get("Authorization"), " ")
	if len(res) == 2 && res[0] == "Bearer" {
		return res[1], nil
	}
	return "", errors.New("Invalid authorization string")
}
