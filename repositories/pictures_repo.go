package repositories

import (
	"Api-Picture/models"
	"gorm.io/gorm"
)

type PictureRepository struct {
	DB *gorm.DB
}

func NewPictureRepository(db *gorm.DB) *PictureRepository {
	return &PictureRepository{DB: db}
}

func (pr *PictureRepository) GetAll(limit int) ([]models.Pictures, error) {
	var pictures []models.Pictures
	err := pr.DB.Limit(limit).Find(&pictures).Error
	return pictures, err
}

func (pr *PictureRepository) GetById(id string) (models.Pictures, error) {
	var picture models.Pictures
	err := pr.DB.Where("id=?", id).First(&picture).Error
	return picture, err
}

func (pr *PictureRepository) Count() (int64, error) {
	var count int64
	err := pr.DB.Model(&models.Pictures{}).Count(&count).Error
	return count, err
}

func (pr *PictureRepository) GetPicturesAfterID(lastSeenID int, limit int) ([]models.Pictures, error) {
	var pictures []models.Pictures
	err := pr.DB.Limit(limit).Offset(lastSeenID).Find(&pictures).Error
	return pictures, err
}

func (pr *PictureRepository) GetPicturesPaginated(lastSeenID int, limit int) ([]models.Pictures, error) {
	var pictures []models.Pictures
	err := pr.DB.Order("id").Limit(limit).Where("id > ?", lastSeenID).Find(&pictures).Error
	return pictures, err
}
