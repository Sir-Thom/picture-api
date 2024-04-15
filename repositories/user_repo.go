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
	err := ur.DB.Create(&user).Error
	return err
}

func (ur *UserRepository) SignIn(user models.User) (models.User, error) {
	var u models.User
	err := ur.DB.Where("email = ? AND password = ?", user.Email, user.Password).First(&u).Error
	return u, err
}
