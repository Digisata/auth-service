package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"golang.org/x/crypto/bcrypt"
)

type ProfileUsecase struct {
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewProfileUsecase(userRepository UserRepository, timeout time.Duration) *ProfileUsecase {
	return &ProfileUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (pu ProfileUsecase) CreateProfile(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	_, err := pu.userRepository.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("user already exists with the given email")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user.Password = string(encryptedPassword)

	return pu.userRepository.Create(ctx, user)
}

func (pu ProfileUsecase) GetProfileByID(c context.Context, userID string) (*domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	user, err := pu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.Profile{Name: user.Name, Email: user.Email}, nil
}
