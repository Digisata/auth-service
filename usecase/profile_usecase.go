package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	memcachedRepo "github.com/digisata/auth-service/repository/memcached"
	mongoRepo "github.com/digisata/auth-service/repository/mongo"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileUsecase struct {
	jwt     *jwtio.JSONWebToken
	cfg     *bootstrap.Config
	ur      ProfileRepository
	cr      CacheRepository
	timeout time.Duration
}

var _ ProfileRepository = (*mongoRepo.ProfileRepository)(nil)
var _ CacheRepository = (*memcachedRepo.CacheRepository)(nil)

func NewProfileUsecase(jwt *jwtio.JSONWebToken, cfg *bootstrap.Config, ur ProfileRepository, cr CacheRepository, timeout time.Duration) *ProfileUsecase {
	return &ProfileUsecase{
		jwt:     jwt,
		cfg:     cfg,
		ur:      ur,
		cr:      cr,
		timeout: timeout,
	}
}

func (uc ProfileUsecase) GetByID(ctx context.Context, profileID string) (domain.Profile, error) {
	var res domain.Profile
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	res, err := uc.ur.GetByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", profileID))
		}

		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}

func (uc ProfileUsecase) ChangePassword(ctx context.Context, req domain.ChangePasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	claims := ctx.Value("claims")
	profileID := claims.(jwt.MapClaims)["id"].(string)

	user, err := uc.ur.GetByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", profileID))
		}

		return status.Error(codes.Internal, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return status.Error(codes.InvalidArgument, "Incorrect password")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	err = uc.ur.ChangePassword(ctx, profileID, string(encryptedPassword))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
