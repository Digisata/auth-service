package usecase

import (
	"context"
	"time"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/internal/tokenutil"
)

type LoginUsecase struct {
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(userRepository UserRepository, timeout time.Duration) *LoginUsecase {
	return &LoginUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu LoginUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, lu.contextTimeout)
	defer cancel()
	return lu.userRepository.GetByEmail(ctx, email)
}

func (lu LoginUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return tokenutil.CreateAccessToken(user, secret, expiry)
}

func (lu LoginUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return tokenutil.CreateRefreshToken(user, secret, expiry)
}
