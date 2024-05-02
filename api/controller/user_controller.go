package controller

import (
	"context"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	userPb "github.com/digisata/auth-service/stubs/user"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	userPb.UnimplementedAuthServiceServer
	LoginUsecase        LoginUsecase
	RefreshTokenUsecase RefreshTokenUsecase
	UserUsecase         UserUsecase
	Env                 *bootstrap.Config
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

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword())) != nil {
		return nil, err
	}

	accessToken, err := uc.LoginUsecase.CreateAccessToken(&user, uc.Env.AccessTokenSecret, uc.Env.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.LoginUsecase.CreateRefreshToken(&user, uc.Env.RefreshTokenSecret, uc.Env.RefreshTokenExpiryHour)
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
	id, err := uc.RefreshTokenUsecase.ExtractIDFromToken(req.GetRefreshToken(), uc.Env.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}

	user, err := uc.RefreshTokenUsecase.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.RefreshTokenUsecase.CreateAccessToken(&user, uc.Env.AccessTokenSecret, uc.Env.AccessTokenExpiryHour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.RefreshTokenUsecase.CreateRefreshToken(&user, uc.Env.RefreshTokenSecret, uc.Env.RefreshTokenExpiryHour)
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
	// userID := c.GetString("x-user-id")

	// user, err := uc.UserUsecase.GetUserByID(c, userID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, user)
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByID not implemented")
}
