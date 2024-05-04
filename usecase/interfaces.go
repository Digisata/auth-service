package usecase

import (
	"context"

	"github.com/digisata/auth-service/domain"
)

type (
	UserRepository interface {
		Create(ctx context.Context, req domain.User) error
		Fetch(ctx context.Context) ([]domain.User, error)
		GetByEmail(ctx context.Context, email string) (domain.User, error)
		GetByID(ctx context.Context, id string) (domain.User, error)
	}
)
