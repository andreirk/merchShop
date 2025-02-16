package auth

import (
	"errors"
	"go/avito-test/internal/models"
	"go/avito-test/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

const (
	ErrPasswordTooShort = "password is too short"
)

type AuthService struct {
	userRepo di.IUserRepository
}

func NewAuthService(userRepo di.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (service *AuthService) Register(username, password string) (string, error) {
	existingUser, _ := service.userRepo.FindUserByName(username)
	if existingUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hasheDPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Password: string(hasheDPass),
		Username: username,
	}
	_, err = service.userRepo.Create(user)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func (service *AuthService) Login(username, password string) (string, error) {
	user, err := service.userRepo.FindUserByName(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return user.Username, nil
}
