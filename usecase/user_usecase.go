package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	memcachedRepo "github.com/digisata/auth-service/repository/memcached"
	mongoRepo "github.com/digisata/auth-service/repository/mongo"
	"go.mongodb.org/mongo-driver/mongo"
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

var _ UserRepository = (*mongoRepo.UserRepository)(nil)
var _ CacheRepository = (*memcachedRepo.CacheRepository)(nil)

func NewUserUsecase(jwt *jwtio.JSONWebToken, cfg *bootstrap.Config, ur UserRepository, cr CacheRepository, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		jwt:     jwt,
		cfg:     cfg,
		ur:      ur,
		cr:      cr,
		timeout: timeout,
	}
}

func (uc UserUsecase) generateToken(user domain.User) (domain.AuthResponse, error) {
	var res domain.AuthResponse
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
		Exp: int32(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
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
		Exp: int32(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	res = domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uc UserUsecase) LoginAdmin(ctx context.Context, req domain.User) (domain.AuthResponse, error) {
	var res domain.AuthResponse
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

	if !user.IsActive || user.DeletedAt != 0 {
		return res, status.Error(codes.Unauthenticated, "Your account has been deleted")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) LoginCustomer(ctx context.Context, req domain.User) (domain.AuthResponse, error) {
	var res domain.AuthResponse
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

	if !user.IsActive || user.DeletedAt != 0 {
		return res, status.Error(codes.Unauthenticated, "Your account has been deleted")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) LoginCommittee(ctx context.Context, req domain.User) (domain.AuthResponse, error) {
	var res domain.AuthResponse
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

	if !user.IsActive || user.DeletedAt != 0 {
		return res, status.Error(codes.Unauthenticated, "Your account has been deleted")
	}

	res, err = uc.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uc UserUsecase) RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.AuthResponse, error) {
	var res domain.AuthResponse
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	claims, err := uc.jwt.VerifyRefreshToken(req.RefreshToken, uc.cfg.Jwt.RefreshTokenSecret)
	if err != nil {
		return res, err
	}

	user, err := uc.ur.GetByID(ctx, claims["id"].(string))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", claims["id"].(string)))
		}

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
		Exp: int32(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
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
		Exp: int32(now.Add(time.Hour * time.Duration(uc.cfg.Jwt.RefreshTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(req.AccessToken)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(req.RefreshToken)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return res, status.Error(codes.Internal, err.Error())
	}

	res = domain.AuthResponse{
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", userID))
		}

		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (uc UserUsecase) Logout(ctx context.Context, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	accessToken, _ := uc.jwt.GetAccessToken(ctx)

	err := uc.cr.Delete(accessToken)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return status.Error(codes.Internal, err.Error())
	}

	err = uc.cr.Delete(refreshToken)
	if err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (uc UserUsecase) Update(ctx context.Context, req domain.UpdateUser) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	userID := req.ID.Hex()

	_, err := uc.ur.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", userID))
		}

		return status.Error(codes.Internal, err.Error())
	}

	err = uc.ur.Update(ctx, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (uc UserUsecase) Delete(ctx context.Context, req domain.DeleteUser) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	userID := req.ID.Hex()

	_, err := uc.ur.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", userID))
		}

		return status.Error(codes.Internal, err.Error())
	}

	err = uc.ur.Delete(ctx, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
