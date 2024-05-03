package usecase

import (
	"context"
	"time"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
)

type LoginUsecase struct {
	jwt            *jwtio.JSONWebToken
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewLoginUsecase(jwt *jwtio.JSONWebToken, userRepository UserRepository, timeout time.Duration) *LoginUsecase {
	return &LoginUsecase{
		jwt:            jwt,
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
	return lu.jwt.CreateAccessToken(user, secret, expiry)
}

func (lu LoginUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return lu.jwt.CreateRefreshToken(user, secret, expiry)
}
