package jwt

import (
	"errors"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("mini-issue-secret")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	gojwt.RegisteredClaims
}


func GenerateToken(userid int64, username string) (string, error) {
	claims := Claims {
		UserID: userid,
		Username: username,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: gojwt.NewNumericDate(time.Now()),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenstring string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenstring, &Claims{}, func(token *gojwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
