package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	jwt            *jwtio.JSONWebToken
	cfg            *bootstrap.Config
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(jwt *jwtio.JSONWebToken, cfg *bootstrap.Config, userRepository UserRepository, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		jwt:            jwt,
		cfg:            cfg,
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu UserUsecase) Login(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, err
	}

	payload := jwtio.Payload{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	accessToken, err := uu.jwt.CreateAccessToken(payload, uu.cfg.Jwt.AccessTokenSecret, uu.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	refreshToken, err := uu.jwt.CreateRefreshToken(payload, uu.cfg.Jwt.RefreshTokenSecret, uu.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	res = domain.Login{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uu UserUsecase) RefreshToken(ctx context.Context, token string) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	id, err := uu.jwt.ExtractIDFromToken(token, uu.cfg.Jwt.RefreshTokenSecret)
	if err != nil {
		return res, err
	}

	user, err := uu.userRepository.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	payload := jwtio.Payload{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	accessToken, err := uu.jwt.CreateAccessToken(payload, uu.cfg.Jwt.AccessTokenSecret, uu.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	refreshToken, err := uu.jwt.CreateRefreshToken(payload, uu.cfg.Jwt.RefreshTokenSecret, uu.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	res = domain.Login{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uu UserUsecase) CreateUser(ctx context.Context, req domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	_, err := uu.userRepository.GetByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("user already exists with the given email")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	req.Password = string(encryptedPassword)

	return uu.userRepository.Create(ctx, req)
}

func (uu UserUsecase) GetUserByID(ctx context.Context, userID string) (domain.UserProfile, error) {
	var res domain.UserProfile
	ctx, cancel := context.WithTimeout(ctx, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return res, err
	}

	res = domain.UserProfile{
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}
