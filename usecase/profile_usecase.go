package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/domain"
	"github.com/digisata/auth-service/pkg/jwtio"
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
	fmt.Println("tc")
	ctx, cancel := context.WithTimeout(ctx, uu.timeout)
	defer cancel()

	res, err := uu.ur.GetByID(ctx, profileID)
	if err != nil {
		return res, status.Error(codes.Internal, err.Error())
	}

	return res, nil
}
