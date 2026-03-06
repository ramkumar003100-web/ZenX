package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string   `json:"uid"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

func (m *JWTManager) Sign(userID string, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
