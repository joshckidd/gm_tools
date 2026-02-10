package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userName string, tokenSecret string, expiresIn time.Duration) (string, error) {
	res := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "dm_tools",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userName,
	})
	return res.SignedString([]byte(tokenSecret))
}

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

func GetBearerToken(headers http.Header) (string, error) {
	res := strings.Split(headers.Get("Authorization"), " ")
	fmt.Println(res)
	if len(res) == 2 && res[0] == "Bearer" {
		return res[1], nil
	}
	return "", errors.New("Invalid authorization string")
}
