package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	UserUsecase interface {
		Login(ctx context.Context, req domain.User) (domain.Login, error)
		RefreshToken(ctx context.Context, token string) (domain.Login, error)
		CreateUser(ctx context.Context, req domain.User) error
		GetUserByID(ctx context.Context, userID string) (domain.UserProfile, error)
		Logout(ctx context.Context, refreshToken string) error
	}
)
