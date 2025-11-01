package services

import (
	"gofiber-crud/models"

	"gorm.io/gorm"
)

type UserService interface {
	FindAll() ([]models.User, error)
	FindById(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteById(id int) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) FindAll() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *userService) FindById(id int) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	return &user, result.Error
}

func (s *userService) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (s *userService) CreateUser(user *models.User) error {
	result := s.db.Create(user)
	return result.Error
}

func (s *userService) UpdateUser(user *models.User) error {
	result := s.db.Save(user)
	return result.Error
}

func (s *userService) DeleteById(id int) error {
	result := s.db.Delete(&models.User{}, id)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
