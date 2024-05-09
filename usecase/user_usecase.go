package usecase

import (
	"context"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (uc UserUsecase) generateToken(user domain.User) (domain.Login, error) {
	var res domain.Login
	payload := jwtio.Payload{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	now := time.Now()

	accessToken, err := uc.jwt.CreateAccessToken(payload, uc.cfg.Jwt.AccessTokenSecret, now, uc.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uc.cr.Set(domain.CacheItem{
		Key: accessToken,
		Exp: int(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := uc.jwt.CreateRefreshToken(payload, uc.cfg.Jwt.RefreshTokenSecret, now, uc.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uc.cr.Set(domain.CacheItem{
		Key: refreshToken,
		Exp: int(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	res = domain.Login{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uc UserUsecase) LoginAdmin(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	user, err := uc.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.ADMIN) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) LoginCustomer(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	user, err := uc.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.CUSTOMER) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) LoginCommittee(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	user, err := uc.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.COMMITTEE) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	claims, err := uc.jwt.VerifyRefreshToken(req.RefreshToken, uc.cfg.Jwt.RefreshTokenSecret)
	if err != nil {
		return res, err
	}

	user, err := uc.ur.GetByID(ctx, claims["id"].(string))
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	payload := jwtio.Payload{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	now := time.Now()

	newAccessToken, err := uc.jwt.CreateAccessToken(payload, uc.cfg.Jwt.AccessTokenSecret, now, uc.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uc.cr.Set(domain.CacheItem{
		Key: newAccessToken,
		Exp: int(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	newRefreshToken, err := uc.jwt.CreateRefreshToken(payload, uc.cfg.Jwt.RefreshTokenSecret, now, uc.cfg.Jwt.RefreshTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uc.cr.Set(domain.CacheItem{
		Key: newRefreshToken,
		Exp: int(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(req.AccessToken)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(req.RefreshToken)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	res = domain.Login{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return res, nil
}

func (uc UserUsecase) Create(ctx context.Context, req domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	_, err := uc.ur.GetByEmail(ctx, req.Email)
	if err == nil {
		return status.Error(codes.InvalidArgument, "User already exists with the given email")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	req.Password = string(encryptedPassword)
	err = uc.ur.Create(ctx, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (uc UserUsecase) GetByID(ctx context.Context, userID string) (domain.User, error) {
	var res domain.User
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	res, err := uc.ur.GetByID(ctx, userID)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (uc UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	accessToken, _ := uc.jwt.GetAccessToken(ctx)

	err := uc.cr.Delete(accessToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(refreshToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
