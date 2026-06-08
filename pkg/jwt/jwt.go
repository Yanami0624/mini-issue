package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("mini-issue-secret")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userid int64, username string) (string, error) {
	claims := Claims{
		UserID:   userid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenstring string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(token *jwt.Token) (any, error) {
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
