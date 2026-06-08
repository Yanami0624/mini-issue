package jwt

import (
	"testing"
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndParseToken(t *testing.T) {
	tokenString, err := GenerateToken(42, "alice")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != 42 {
		t.Fatalf("claims.UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "alice" {
		t.Fatalf("claims.Username = %q, want alice", claims.Username)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatalf("claims.ExpiresAt = %v, want future expiration", claims.ExpiresAt)
	}
}

func TestParseTokenRejectsInvalidToken(t *testing.T) {
	if _, err := ParseToken("not-a-token"); err == nil {
		t.Fatal("ParseToken() error = nil, want error")
	}
}

func TestParseTokenRejectsExpiredToken(t *testing.T) {
	claims := Claims{
		UserID:   42,
		Username: "alice",
		RegisteredClaims: golangjwt.RegisteredClaims{
			ExpiresAt: golangjwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  golangjwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	tokenString, err := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := ParseToken(tokenString); err == nil {
		t.Fatal("ParseToken() error = nil, want error")
	}
}
