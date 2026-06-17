package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour

	TypeAccess  = "access"
	TypeRefresh = "refresh"
)

var (
	ErrEmptySecretKey = errors.New("jwt secret key must not be empty")
	ErrInvalidToken   = errors.New("invalid token")
	ErrWrongTokenType = errors.New("wrong token type")
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId uint, email string, name string, role string) (string, error)
	GenerateRefreshToken(userId uint) (string, error)
	ValidateToken(tokenStr string) (*JWTClaims, error)
	ValidateRefreshToken(tokenStr string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTService(secretKey string) (JWTService, error) {
	if secretKey == "" {
		return nil, ErrEmptySecretKey
	}

	return &jwtService{
		secretKey:            secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}, nil
}

func (js *jwtService) generate(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(js.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (js *jwtService) GenerateToken(userId uint, email string, name string, role string) (string, error) {
	claims := JWTClaims{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
		Type:   TypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-library",
		},
	}
	return js.generate(claims)
}

func (js *jwtService) GenerateRefreshToken(userId uint) (string, error) {
	claims := JWTClaims{
		UserID: userId,
		Type:   TypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-library",
		},
	}
	return js.generate(claims)
}

func (js *jwtService) parse(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (js *jwtService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	claims, err := js.parse(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != TypeAccess {
		return nil, ErrWrongTokenType
	}
	return claims, nil
}

func (js *jwtService) ValidateRefreshToken(tokenStr string) (*JWTClaims, error) {
	claims, err := js.parse(tokenStr)
	if err != nil {
		return nil, err
	}
	if claims.Type != TypeRefresh {
		return nil, ErrWrongTokenType
	}
	return claims, nil
}
