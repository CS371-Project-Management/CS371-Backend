package services

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
	}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uint) (models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) CreateUser(user *models.User) error {
	// hashPassword(user)
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
