package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	LoginUsecase interface {
		GetUserByEmail(c context.Context, email string) (domain.User, error)
		CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error)
		CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error)
	}

	RefreshTokenUsecase interface {
		GetUserByID(c context.Context, id string) (domain.User, error)
		CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error)
		CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error)
		ExtractIDFromToken(requestToken string, secret string) (string, error)
	}

	UserUsecase interface {
		CreateUser(c context.Context, user *domain.User) error
		GetUserByID(c context.Context, userID string) (*domain.UserProfile, error)
	}
)
