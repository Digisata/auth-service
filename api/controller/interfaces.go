package controller

import (
	"context"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
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

	ProfileUsecase interface {
		CreateProfile(c context.Context, user *domain.User) error
		GetProfileByID(c context.Context, userID string) (*domain.Profile, error)
	}
)
