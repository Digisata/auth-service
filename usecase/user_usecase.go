package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/digisata/auth-service/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepository UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepository UserRepository, timeout time.Duration) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu UserUsecase) CreateUser(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	_, err := uu.userRepository.GetByEmail(ctx, user.Email)
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

	return uu.userRepository.Create(ctx, user)
}

func (uu UserUsecase) GetUserByID(c context.Context, userID string) (*domain.UserProfile, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserProfile{Name: user.Name, Email: user.Email}, nil
}
