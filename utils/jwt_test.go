package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SisyphusSQ/go-starter/v2/config"
)

func TestTokenRoundTrip(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "0123456789abcdef0123456789abcdef",
		Expire: time.Hour,
		Issuer: "test-service",
	}
	token, err := GenerateToken(42, "user@example.com", cfg)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	claims, err := ParseToken(token, cfg.Secret, cfg.Issuer)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 42 || claims.Email != "user@example.com" {
		t.Fatalf("claims = %#v", claims)
	}
}

func TestParseTokenRejectsWrongAlgorithm(t *testing.T) {
	claims := &Claims{RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Issuer:    "test-service",
	}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).
		SignedString([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	if _, err = ParseToken(token, "0123456789abcdef0123456789abcdef", "test-service"); err == nil {
		t.Fatal("ParseToken() accepted HS384")
	}
}

func TestTokenRotationAndClaims(t *testing.T) {
	cfg := config.JWTConfig{Secret: UUID(), Expire: time.Hour, Issuer: "test-service"}
	first, err := GenerateToken(42, "user@example.com", cfg)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateToken(42, "user@example.com", cfg)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("new session reused previous token")
	}
	if _, err := ParseToken(first, cfg.Secret, "different-service"); err == nil {
		t.Fatal("wrong issuer accepted")
	}
	for _, claims := range []Claims{
		{RegisteredClaims: jwt.RegisteredClaims{Issuer: cfg.Issuer}},
		{RegisteredClaims: jwt.RegisteredClaims{Issuer: cfg.Issuer, ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))}},
		{RegisteredClaims: jwt.RegisteredClaims{Issuer: cfg.Issuer, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute))}},
	} {
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.Secret))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseToken(token, cfg.Secret, cfg.Issuer); err == nil {
			t.Fatal("invalid temporal claims accepted")
		}
	}
}
