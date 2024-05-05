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
	jwt     *jwtio.JSONWebToken
	cfg     *bootstrap.Config
	ur      UserRepository
	cr      CacheRepository
	timeout time.Duration
}

func NewUserUsecase(jwt *jwtio.JSONWebToken, cfg *bootstrap.Config, ur UserRepository, cr CacheRepository, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		jwt:     jwt,
		cfg:     cfg,
		ur:      ur,
		cr:      cr,
		timeout: timeout,
	}
}

func (uu UserUsecase) Login(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByEmail(ctx, req.Email)
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

	now := time.Now()

	accessToken, err := uu.jwt.CreateAccessToken(payload, uu.cfg.Jwt.AccessTokenSecret, now, uu.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uu.cr.Set(domain.CacheItem{
		Key: accessToken,
		Exp: int(now.Add(time.Hour * time.Duration(uu.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, err
	}

	refreshToken, err := uu.jwt.CreateRefreshToken(payload, uu.cfg.Jwt.RefreshTokenSecret, now, uu.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uu.cr.Set(domain.CacheItem{
		Key: refreshToken,
		Exp: int(now.Add(time.Hour * time.Duration(uu.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, err
	}

	res = domain.Login{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uu UserUsecase) RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	jwt, err := uu.jwt.VerifyRefreshToken(req.RefreshToken, uu.cfg.Jwt.RefreshTokenSecret)
	if err != nil {
		return res, err
	}

	id, err := uu.jwt.ExtractIDFromToken(jwt)
	if err != nil {
		return res, err
	}

	user, err := uu.ur.GetByID(ctx, id)
	if err != nil {
		return res, err
	}

	payload := jwtio.Payload{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	now := time.Now()

	newAccessToken, err := uu.jwt.CreateAccessToken(payload, uu.cfg.Jwt.AccessTokenSecret, now, uu.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uu.cr.Set(domain.CacheItem{
		Key: newAccessToken,
		Exp: int(now.Add(time.Hour * time.Duration(uu.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, err
	}

	newRefreshToken, err := uu.jwt.CreateRefreshToken(payload, uu.cfg.Jwt.RefreshTokenSecret, now, uu.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uu.cr.Set(domain.CacheItem{
		Key: newRefreshToken,
		Exp: int(now.Add(time.Hour * time.Duration(uu.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, err
	}

	err = uu.cr.Delete(req.AccessToken)
	if err != nil {
		return res, err
	}

	err = uu.cr.Delete(req.RefreshToken)
	if err != nil {
		return res, err
	}

	res = domain.Login{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return res, nil
}

func (uu UserUsecase) CreateUser(ctx context.Context, req domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	_, err := uu.ur.GetByEmail(ctx, req.Email)
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

	return uu.ur.Create(ctx, req)
}

func (uu UserUsecase) GetUserByID(ctx context.Context, userID string) (domain.UserProfile, error) {
	var res domain.UserProfile
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByID(ctx, userID)
	if err != nil {
		return res, err
	}

	res = domain.UserProfile{
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}

func (uu UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	accessToken, _ := uu.jwt.GetAccessToken(ctx)

	err := uu.cr.Delete(accessToken)
	if err != nil {
		return err
	}

	err = uu.cr.Delete(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
