package utils

import (
	"errors"
	"strconv"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SisyphusSQ/go-starter/v2/config"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, email string, cfg config.JWTConfig) (string, error) {
	if cfg.Secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	now := time.Now()
	expireAt := now.Add(cfg.Expire)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    cfg.Issuer,
			Subject:   strconv.FormatInt(userID, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenStr, secret, issuer string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("jwt token is empty")
	}
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	if issuer == "" {
		return nil, errors.New("jwt issuer is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(*jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid jwt token")
	}
	return claims, nil
}
