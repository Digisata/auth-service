// Package jwtio is shared pkg of json web token
package jwtio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	FAILED_TO_EXTRACT         = "failed to extract jwt payload"
	UNEXPECTED_SIGNING_METHOD = "unexpected signing method: %v"
)

type JSONWebToken struct {
	cfg *bootstrap.Config
}

func NewJSONWebToken(cfg *bootstrap.Config) *JSONWebToken {
	return &JSONWebToken{
		cfg: cfg,
	}
}

func (j JSONWebToken) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &domain.JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (j JSONWebToken) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return rt, nil
}

func (j JSONWebToken) Validate(ctx context.Context) (jwt.MapClaims, error) {
	accessToken, err := getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return j.validateLoginToken(token)
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf(FAILED_TO_EXTRACT)
	}

	return claims, nil
}

func (j JSONWebToken) validateLoginToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf(UNEXPECTED_SIGNING_METHOD, token.Header["alg"])
	}

	return []byte(j.cfg.AccessTokenSecret), nil
}

func getAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	split := strings.Split(values[0], " ")
	if len(split) != 2 {
		return "", status.Error(codes.Unauthenticated, "invalid access token format")
	}

	if split[0] != "Bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid access token format")
	}

	return split[1], nil
}

func (j *JSONWebToken) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(UNEXPECTED_SIGNING_METHOD, token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", fmt.Errorf(FAILED_TO_EXTRACT)
	}

	return claims["id"].(string), nil
}
