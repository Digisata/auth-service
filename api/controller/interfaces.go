package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	UserUsecase interface {
		LoginAdmin(ctx context.Context, req domain.User) (domain.Login, error)
		LoginCustomer(ctx context.Context, req domain.User) (domain.Login, error)
		LoginCommittee(ctx context.Context, req domain.User) (domain.Login, error)
		RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.Login, error)
		Create(ctx context.Context, req domain.User) error
		GetByID(ctx context.Context, id string) (domain.User, error)
		Logout(ctx context.Context, refreshToken string) error
	}

	ProfileUsecase interface {
		GetByID(ctx context.Context, id string) (domain.Profile, error)
	}
)
