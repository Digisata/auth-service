package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	UserUsecase interface {
		LoginAdmin(ctx context.Context, req domain.User) (domain.AuthResponse, error)
		LoginCustomer(ctx context.Context, req domain.User) (domain.AuthResponse, error)
		LoginCommittee(ctx context.Context, req domain.User) (domain.AuthResponse, error)
		RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.AuthResponse, error)
		Create(ctx context.Context, req domain.User) error
		GetAll(ctx context.Context, req domain.GetAllUserRequest) ([]domain.User, error)
		GetByID(ctx context.Context, id string) (domain.User, error)
		Update(ctx context.Context, req domain.UpdateUser) error
		Delete(ctx context.Context, req domain.DeleteUser) error
		Logout(ctx context.Context, refreshToken string) error
	}

	ProfileUsecase interface {
		GetByID(ctx context.Context, id string) (domain.Profile, error)
		ChangePassword(ctx context.Context, req domain.ChangePasswordRequest) error
	}
)
