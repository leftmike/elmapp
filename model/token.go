package model

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	signingKey      = []byte("TOPSECRETSIGNINGKEY")
	errInvalidToken = errors.New("invalid token")
)

type tokenClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func newToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		tokenClaims{
			Username: user.Username,
			Email:    user.Email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				Issuer:    "elmapp",
			},
		})
	return token.SignedString(signingKey)
}

func ValidateToken(s string) (*User, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	}

	token, err := jwt.ParseWithClaims(s, &tokenClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !token.Valid || !ok {
		return nil, errInvalidToken
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	user, ok := userByUsername[claims.Username]
	if !ok || user.Email != claims.Email {
		return nil, errInvalidToken
	}
	return user, nil
}
