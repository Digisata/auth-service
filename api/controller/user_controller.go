package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
	userPb "github.com/digisata/auth-service/stubs/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userPb.UnimplementedAuthServiceServer
	UserUsecase UserUsecase
}

func (uc UserController) LoginAdmin(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := uc.UserUsecase.LoginAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &userPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (uc UserController) LoginCustomer(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := uc.UserUsecase.LoginCustomer(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &userPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (uc UserController) LoginCommittee(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := uc.UserUsecase.LoginCommittee(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &userPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (uc UserController) RefreshToken(ctx context.Context, req *userPb.RefreshTokenRequest) (*userPb.RefreshTokenResponse, error) {
	refreshTokenRequest := domain.RefreshTokenRequest{
		AccessToken:  req.GetAccessToken(),
		RefreshToken: req.GetRefreshToken(),
	}

	data, err := uc.UserUsecase.RefreshToken(ctx, refreshTokenRequest)
	if err != nil {
		return nil, err
	}

	res := &userPb.RefreshTokenResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (uc UserController) CreateUser(ctx context.Context, req *userPb.CreateUserRequest) (*userPb.BaseResponse, error) {
	user := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     int8(req.GetRole()),
	}

	err := uc.UserUsecase.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &userPb.BaseResponse{
		Message: "success",
	}

	return res, nil
}

func (uc UserController) GetUserByID(ctx context.Context, req *empty.Empty) (*userPb.GetUserByIDResponse, error) {
	claims := ctx.Value("claims")

	data, err := uc.UserUsecase.GetUserByID(ctx, claims.(jwt.MapClaims)["id"].(string))
	if err != nil {
		return nil, err
	}

	res := &userPb.GetUserByIDResponse{
		Name:  data.Name,
		Email: data.Email,
	}

	return res, nil
}

func (uc UserController) Logout(ctx context.Context, req *userPb.LogoutRequest) (*userPb.BaseResponse, error) {
	err := uc.UserUsecase.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	res := &userPb.BaseResponse{
		Message: "success",
	}

	return res, nil
}
