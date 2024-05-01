package usecase

import (
	"context"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
)

type (
	UserRepository interface {
		Create(c context.Context, user *domain.User) error
		Fetch(c context.Context) ([]domain.User, error)
		GetByEmail(c context.Context, email string) (domain.User, error)
		GetByID(c context.Context, id string) (domain.User, error)
	}
)
