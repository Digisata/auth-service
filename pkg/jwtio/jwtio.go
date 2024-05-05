// Package jwtio is shared pkg of json web token
package jwtio

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/digisata/auth-service/pkg/constants"
	"github.com/digisata/auth-service/pkg/memcached"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type (
	Config struct {
		AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
		RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
		AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
		RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	}

	Payload struct {
		ID    string
		Name  string
		Email string
		Role  string
	}

	JSONWebToken struct {
		cfg         *Config
		memcachedDB *memcached.Database
	}

	JwtCustomClaims struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	JwtCustomRefreshClaims struct {
		ID string `json:"id"`
		jwt.RegisteredClaims
	}
)

func NewJSONWebToken(cfg *Config, memcachedDB *memcached.Database) *JSONWebToken {
	return &JSONWebToken{
		cfg:         cfg,
		memcachedDB: memcachedDB,
	}
}

func (j JSONWebToken) CreateAccessToken(payload Payload, secret string, now time.Time, expiry int) (accessToken string, err error) {
	claims := &JwtCustomClaims{
		Name: payload.Name,
		ID:   payload.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(expiry))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (j JSONWebToken) CreateRefreshToken(payload Payload, secret string, now time.Time, expiry int) (refreshToken string, err error) {
	claims := &JwtCustomRefreshClaims{
		ID: payload.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(expiry))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (j JSONWebToken) Verify(ctx context.Context) (jwt.MapClaims, error) {
	accessToken, err := j.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	_, err = j.memcachedDB.Get(accessToken)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, status.Error(codes.Unauthenticated, constants.TOKEN_EXPIRED)
		}

		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return j.validateToken(token, j.cfg.AccessTokenSecret)
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, status.Error(codes.Unauthenticated, constants.FAILED_TO_EXTRACT)
	}

	return claims, nil
}

func (j JSONWebToken) VerifyRefreshToken(refreshToken, secret string) (*jwt.Token, error) {
	_, err := j.memcachedDB.Get(refreshToken)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, status.Error(codes.Unauthenticated, constants.REFRESH_TOKEN_EXPIRED)
		}

		return nil, err
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return j.validateToken(token, j.cfg.RefreshTokenSecret)
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j JSONWebToken) ExtractIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", status.Error(codes.Unauthenticated, constants.FAILED_TO_EXTRACT)
	}

	return claims["id"].(string), nil
}

func (j JSONWebToken) validateToken(token *jwt.Token, secret string) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, status.Error(codes.Unauthenticated, constants.UNEXPECTED_SIGNING_METHOD)

	}

	return []byte(secret), nil
}

func (j JSONWebToken) GetAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "Authorization token is not provided")
	}

	split := strings.Split(values[0], " ")
	if len(split) != 2 {
		return "", status.Error(codes.Unauthenticated, "Invalid access token format")
	}

	if split[0] != "Bearer" {
		return "", status.Error(codes.Unauthenticated, "Invalid access token format")
	}

	return split[1], nil
}
