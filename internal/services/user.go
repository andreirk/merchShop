package services

import (
	"go/avito-test/internal/models"
	"go/avito-test/internal/repositories"
)

type IUserService interface {
	FindOrCreateUser(username string) (*models.User, error)
	GetUserByName(username string) (*models.User, error)
	GetUserInfo(username string) (*models.User, error)
}

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserInfo(username string) (*models.User, error) {
	return s.userRepo.FindUserByName(username)
}

func (s *UserService) CreateUser(username string) (*models.User, error) {
	user := &models.User{Username: username}
	if _, err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindOrCreateUser(username string) (*models.User, error) {
	user := &models.User{Username: username}
	if _, err := s.userRepo.FindOrCreate(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) GetUserByName(username string) (*models.User, error) {
	return s.userRepo.FindUserByName(username)
}
