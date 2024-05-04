package controller

import (
	"context"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	userPb "github.com/digisata/auth-service/stubs/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userPb.UnimplementedAuthServiceServer
	Config              *bootstrap.Config
	LoginUsecase        LoginUsecase
	RefreshTokenUsecase RefreshTokenUsecase
	UserUsecase         UserUsecase
}

func (uc UserController) CreateUser(ctx context.Context, req *userPb.CreateUserRequest) (*userPb.BaseResponse, error) {
	user := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	err := uc.UserUsecase.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	res := &userPb.BaseResponse{
		Message: "success",
	}

	return res, nil
}

func (uc UserController) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	user, err := uc.LoginUsecase.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.LoginUsecase.CreateAccessToken(&user, uc.Config.AccessTokenSecret, uc.Config.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.LoginUsecase.CreateRefreshToken(&user, uc.Config.RefreshTokenSecret, uc.Config.RefreshTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	res := &userPb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uc UserController) RefreshToken(ctx context.Context, req *userPb.RefreshTokenRequest) (*userPb.RefreshTokenResponse, error) {
	id, err := uc.RefreshTokenUsecase.ExtractIDFromToken(req.GetRefreshToken(), uc.Config.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}

	user, err := uc.RefreshTokenUsecase.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.RefreshTokenUsecase.CreateAccessToken(&user, uc.Config.AccessTokenSecret, uc.Config.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.RefreshTokenUsecase.CreateRefreshToken(&user, uc.Config.RefreshTokenSecret, uc.Config.RefreshTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	res := &userPb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}

func (uc UserController) GetUserByID(ctx context.Context, req *empty.Empty) (*userPb.GetUserByIDResponse, error) {
	claims := ctx.Value("claims")

	user, err := uc.UserUsecase.GetUserByID(ctx, claims.(jwt.MapClaims)["id"].(string))
	if err != nil {
		return nil, err
	}

	res := &userPb.GetUserByIDResponse{
		Name:  user.Name,
		Email: user.Email,
	}

	return res, nil
}
