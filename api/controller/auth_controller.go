package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
	authPb "github.com/digisata/auth-service/stubs/auth"
	"github.com/digisata/auth-service/usecase"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthController struct {
	authPb.UnimplementedAuthServiceServer
	UserUsecase    UserUsecase
	ProfileUsecase ProfileUsecase
}

var _ UserUsecase = (*usecase.UserUsecase)(nil)
var _ ProfileUsecase = (*usecase.ProfileUsecase)(nil)

// User
func (c AuthController) LoginAdmin(ctx context.Context, req *authPb.LoginRequest) (*authPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &authPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) LoginCustomer(ctx context.Context, req *authPb.LoginRequest) (*authPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginCustomer(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &authPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) LoginCommittee(ctx context.Context, req *authPb.LoginRequest) (*authPb.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginCommittee(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &authPb.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) RefreshToken(ctx context.Context, req *authPb.RefreshTokenRequest) (*authPb.RefreshTokenResponse, error) {
	refreshTokenRequest := domain.RefreshTokenRequest{
		AccessToken:  req.GetAccessToken(),
		RefreshToken: req.GetRefreshToken(),
	}

	data, err := c.UserUsecase.RefreshToken(ctx, refreshTokenRequest)
	if err != nil {
		return nil, err
	}

	res := &authPb.RefreshTokenResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) CreateUser(ctx context.Context, req *authPb.CreateUserRequest) (*authPb.BaseResponse, error) {
	user := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     int8(req.GetRole()),
	}

	err := c.UserUsecase.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &authPb.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (c AuthController) GetUserByID(ctx context.Context, req *authPb.GetUserByIDRequest) (*authPb.GetUserByIDResponse, error) {
	data, err := c.UserUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := &authPb.GetUserByIDResponse{
		Id:        data.ID.Hex(),
		Name:      data.Name,
		Email:     data.Email,
		Role:      int32(data.Role),
		IsActive:  data.IsActive,
		Note:      data.Note,
		CreatedAt: int32(data.CreatedAt),
		UpdatedAt: int32(data.UpdatedAt),
		DeletedAt: int32(data.DeletedAt),
	}

	return res, nil
}

func (c AuthController) Logout(ctx context.Context, req *authPb.LogoutRequest) (*authPb.BaseResponse, error) {
	err := c.UserUsecase.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	res := &authPb.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

// Profile
func (c AuthController) GetProfileByID(ctx context.Context, req *emptypb.Empty) (*authPb.GetProfileByIDResponse, error) {
	claims := ctx.Value("claims")

	data, err := c.ProfileUsecase.GetByID(ctx, claims.(jwt.MapClaims)["id"].(string))
	if err != nil {
		return nil, err
	}

	res := &authPb.GetProfileByIDResponse{
		Id:        data.ID.Hex(),
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: int32(data.CreatedAt),
		UpdatedAt: int32(data.UpdatedAt),
		DeletedAt: int32(data.DeletedAt),
	}

	return res, nil
}

func (c AuthController) ChangePassword(ctx context.Context, req *authPb.ChangePasswordRequest) (*authPb.BaseResponse, error) {
	changePasswordRequest := domain.ChangePasswordRequest{
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}

	err := c.ProfileUsecase.ChangePassword(ctx, changePasswordRequest)
	if err != nil {
		return nil, err
	}

	res := &authPb.BaseResponse{
		Message: "Success",
	}

	return res, nil
}
