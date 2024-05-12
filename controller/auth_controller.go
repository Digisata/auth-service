package controller

import (
	"context"

	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/stubs"
	"github.com/digisata/auth-service/usecase"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthController struct {
	stubs.UnimplementedAuthServiceServer
	UserUsecase    UserUsecase
	ProfileUsecase ProfileUsecase
}

var _ UserUsecase = (*usecase.UserUsecase)(nil)
var _ ProfileUsecase = (*usecase.ProfileUsecase)(nil)

// User
func (c AuthController) LoginAdmin(ctx context.Context, req *stubs.LoginRequest) (*stubs.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginAdmin(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &stubs.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) LoginCustomer(ctx context.Context, req *stubs.LoginRequest) (*stubs.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginCustomer(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &stubs.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) LoginCommittee(ctx context.Context, req *stubs.LoginRequest) (*stubs.LoginResponse, error) {
	payload := domain.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	data, err := c.UserUsecase.LoginCommittee(ctx, payload)
	if err != nil {
		return nil, err
	}

	res := &stubs.LoginResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) RefreshToken(ctx context.Context, req *stubs.RefreshTokenRequest) (*stubs.RefreshTokenResponse, error) {
	refreshTokenRequest := domain.RefreshTokenRequest{
		AccessToken:  req.GetAccessToken(),
		RefreshToken: req.GetRefreshToken(),
	}

	data, err := c.UserUsecase.RefreshToken(ctx, refreshTokenRequest)
	if err != nil {
		return nil, err
	}

	res := &stubs.RefreshTokenResponse{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}

	return res, nil
}

func (c AuthController) CreateUser(ctx context.Context, req *stubs.CreateUserRequest) (*stubs.BaseResponse, error) {
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

	res := &stubs.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (c AuthController) GetAllUser(ctx context.Context, req *stubs.GetAllUserRequest) (*stubs.GetAllUserResponse, error) {
	filter := domain.GetAllUserRequest{
		Search:   req.GetSearch(),
		IsActive: req.GetIsActive(),
	}

	users, err := c.UserUsecase.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	res := &stubs.GetAllUserResponse{}
	for _, user := range users {
		data := &stubs.GetUserByIDResponse{
			Id:        user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      int32(user.Role),
			IsActive:  user.IsActive,
			Note:      user.Note,
			CreatedAt: int32(user.CreatedAt),
			UpdatedAt: int32(user.UpdatedAt),
			DeletedAt: int32(user.DeletedAt),
		}

		res.Users = append(res.Users, data)
	}

	return res, nil
}
func (c AuthController) GetUserByID(ctx context.Context, req *stubs.GetUserByIDRequest) (*stubs.GetUserByIDResponse, error) {
	data, err := c.UserUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := &stubs.GetUserByIDResponse{
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

func (c AuthController) UpdateUser(ctx context.Context, req *stubs.UpdateUserRequest) (*stubs.BaseResponse, error) {
	idHex, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := domain.UpdateUser{
		ID:       idHex,
		Name:     req.GetName(),
		IsActive: req.GetIsActive(),
		Note:     req.GetNote(),
	}

	err = c.UserUsecase.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &stubs.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (c AuthController) DeleteUser(ctx context.Context, req *stubs.DeleteUserRequest) (*stubs.BaseResponse, error) {
	idHex, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := domain.DeleteUser{
		ID: idHex,
	}

	err = c.UserUsecase.Delete(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &stubs.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (c AuthController) Logout(ctx context.Context, req *stubs.LogoutRequest) (*stubs.BaseResponse, error) {
	err := c.UserUsecase.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	res := &stubs.BaseResponse{
		Message: "Success",
	}

	return res, nil
}

// Profile
func (c AuthController) GetProfileByID(ctx context.Context, req *emptypb.Empty) (*stubs.GetProfileByIDResponse, error) {
	claims := ctx.Value("claims")

	data, err := c.ProfileUsecase.GetByID(ctx, claims.(jwt.MapClaims)["id"].(string))
	if err != nil {
		return nil, err
	}

	res := &stubs.GetProfileByIDResponse{
		Id:        data.ID.Hex(),
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: int32(data.CreatedAt),
		UpdatedAt: int32(data.UpdatedAt),
		DeletedAt: int32(data.DeletedAt),
	}

	return res, nil
}

func (c AuthController) ChangePassword(ctx context.Context, req *stubs.ChangePasswordRequest) (*stubs.BaseResponse, error) {
	changePasswordRequest := domain.ChangePasswordRequest{
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	}

	err := c.ProfileUsecase.ChangePassword(ctx, changePasswordRequest)
	if err != nil {
		return nil, err
	}

	res := &stubs.BaseResponse{
		Message: "Success",
	}

	return res, nil
}
