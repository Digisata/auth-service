package usecase

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	UserRepository interface {
		Create(ctx context.Context, req domain.User) error
		GetAll(ctx context.Context, req domain.GetAllUserRequest) ([]domain.User, error)
		GetByEmail(ctx context.Context, email string) (domain.User, error)
		GetByID(ctx context.Context, id string) (domain.User, error)
		Update(ctx context.Context, req domain.UpdateUser) error
		Delete(ctx context.Context, req domain.DeleteUser) error
	}

	ProfileRepository interface {
		GetByID(ctx context.Context, id string) (domain.Profile, error)
		ChangePassword(ctx context.Context, id, newPassword string) error
	}

	CacheRepository interface {
		Set(req domain.CacheItem) error
		Get(key string) (domain.CacheItem, error)
		Delete(key string) error
	}
)
