package service

import (
	"github.com/Gursevak56/food-delivery-platform/services/user-service/models"
	"github.com/Gursevak56/food-delivery-platform/services/user-service/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) CreateUser(user *models.User) (string, error) {
	return s.Repo.CreateUser(user)
}
func (s *UserService) LoginUser(credentials models.LoginCredentials) (string, error) {
	return s.Repo.LoginUser(credentials)
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.Repo.DeleteUser(id)
}
