package usecase

import (
	"context"
	"time"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
)

type RefreshTokenUsecase struct {
	jwt            *jwtio.JSONWebToken
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewRefreshTokenUsecase(jwt *jwtio.JSONWebToken, userRepository UserRepository, timeout time.Duration) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{
		jwt:            jwt,
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (rtu RefreshTokenUsecase) GetUserByID(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, rtu.contextTimeout)
	defer cancel()
	return rtu.userRepository.GetByID(ctx, email)
}

func (rtu RefreshTokenUsecase) CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	return rtu.jwt.CreateAccessToken(user, secret, expiry)
}

func (rtu RefreshTokenUsecase) CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	return rtu.jwt.CreateRefreshToken(user, secret, expiry)
}

func (rtu RefreshTokenUsecase) ExtractIDFromToken(requestToken string, secret string) (string, error) {
	return rtu.jwt.ExtractIDFromToken(requestToken, secret)
}
