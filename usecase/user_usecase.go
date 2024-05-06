package usecase

import (
	"context"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	"github.com/golang-jwt/jwt/v4"
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

func (uu UserUsecase) generateToken(user domain.User) (domain.Login, error) {
	var res domain.Login
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
		return res, status.Error(codes.Internal, err.Error())
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
		return res, status.Error(codes.Internal, err.Error())
	}

	res = domain.Login{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uu UserUsecase) LoginAdmin(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.Admin) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uu.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uu UserUsecase) LoginCustomer(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.Customer) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uu.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uu UserUsecase) LoginCommittee(ctx context.Context, req domain.User) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	if user.Role != int8(domain.Committee) {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, status.Error(codes.InvalidArgument, "Incorrect email or password")
	}

	res, err = uu.generateToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (uu UserUsecase) RefreshToken(ctx context.Context, req domain.RefreshTokenRequest) (domain.Login, error) {
	var res domain.Login
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	claims, err := uu.jwt.VerifyRefreshToken(req.RefreshToken, uu.cfg.Jwt.RefreshTokenSecret)
	if err != nil {
		return res, err
	}

	user, err := uu.ur.GetByID(ctx, claims["id"].(string))
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

	newAccessToken, err := uu.jwt.CreateAccessToken(payload, uu.cfg.Jwt.AccessTokenSecret, now, uu.cfg.Jwt.AccessTokenExpiryHour)
	if err != nil {
		return res, err
	}

	err = uu.cr.Set(domain.CacheItem{
		Key: newAccessToken,
		Exp: int(now.Add(time.Hour * time.Duration(uu.cfg.Jwt.AccessTokenExpiryHour)).Unix()),
	})
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
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
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uu.cr.Delete(req.AccessToken)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	err = uu.cr.Delete(req.RefreshToken)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
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

	claims := ctx.Value("claims")
	role := claims.(jwt.MapClaims)["role"].(int8)

	if role != int8(domain.Admin) {
		return status.Error(codes.Unauthenticated, "Only admin allowed")
	}

	_, err := uu.ur.GetByEmail(ctx, req.Email)
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
	err = uu.ur.Create(ctx, req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (uu UserUsecase) GetUserByID(ctx context.Context, userID string) (domain.UserProfile, error) {
	var res domain.UserProfile
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	user, err := uu.ur.GetByID(ctx, userID)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
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
		return status.Error(codes.Internal, err.Error())
	}

	err = uu.cr.Delete(refreshToken)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
