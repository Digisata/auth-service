package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewProfileUsecase(jwt *jwtio.JSONWebToken, cfg *bootstrap.Config, ur ProfileRepository, cr CacheRepository, timeout time.Duration) *ProfileUsecase {
	return &ProfileUsecase{
		jwt:     jwt,
		cfg:     cfg,
		ur:      ur,
		cr:      cr,
		timeout: timeout,
	}
}

func (uu ProfileUsecase) GetByID(ctx context.Context, profileID string) (domain.Profile, error) {
	var res domain.Profile
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	res, err := uu.ur.GetByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, status.Error(codes.NotFound, fmt.Sprintf("User with id %s not found", profileID))
		}

		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}
