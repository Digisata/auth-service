// Package jwtio is shared pkg of json web token
package jwtio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	FAILED_TO_EXTRACT         = "failed to extract jwt payload"
	UNEXPECTED_SIGNING_METHOD = "unexpected signing method: %v"
)

var claims jwt.MapClaims

type JSONWebToken struct {
	cfg *bootstrap.Config
}

func NewJSONWebToken(cfg *bootstrap.Config) *JSONWebToken {
	return &JSONWebToken{
		cfg: cfg,
	}
}

func (j *JSONWebToken) Generate(ctx context.Context, payload interface{}, privateKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)
	claims["dat"] = payload    // Our custom data.
	claims["iat"] = now.Unix() // The time at which the token was issued.

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JSONWebToken) Validate(ctx context.Context) (jwt.MapClaims, error) {
	accessToken, err := getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return j.validateLoginToken(ctx, token, &claims)
	})
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (j *JSONWebToken) validateLoginToken(ctx context.Context, token *jwt.Token, claim *jwt.MapClaims) (interface{}, error) {
	c, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf(FAILED_TO_EXTRACT)
	}

	*claim = c

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

func (j *JSONWebToken) GetJWTClaim(ctx context.Context) (jwt.MapClaims, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md["authorization"]
	split := strings.Split(values[0], " ")
	accessToken := split[1]

	claims := jwt.MapClaims{}

	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		c, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf(FAILED_TO_EXTRACT)
		}

		claims = c

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(UNEXPECTED_SIGNING_METHOD, token.Header["alg"])
		}

		return []byte(j.cfg.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
