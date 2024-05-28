package repositories

import (
	"Api-Picture/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) SignUp(user models.User) error {
	return ur.DB.Create(&user).Error
}

func (ur *UserRepository) Login(email string) (models.User, error) {
	var user models.User
	err := ur.DB.Where("email = ?", email).First(&user).Error
	return user, err
}
