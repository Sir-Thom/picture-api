package repositories

import (
	"Api-Picture/models"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"time"
)

type PictureRepository struct {
	DB    *gorm.DB
	Cache *cache.Cache
}

func NewPictureRepository(db *gorm.DB) *PictureRepository {
	return &PictureRepository{DB: db, Cache: cache.New(5*time.Minute, 10*time.Minute)}
}

func (pr *PictureRepository) GetAll(limit int) ([]models.Pictures, error) {

	var pictures []models.Pictures
	if pictures, found := pr.Cache.Get("all_pictures"); found {
		return pictures.([]models.Pictures), nil
	}
	err := pr.DB.Select("id, filename, data").Limit(limit).Find(&pictures).Error
	return pictures, err
}

func (pr *PictureRepository) GetById(id string) (models.Pictures, error) {
	var picture models.Pictures
	err := pr.DB.Where("id=?", id).Select("id, filename, data").First(&picture).Error
	return picture, err
}

func (pr *PictureRepository) Count() (int64, error) {
	var count int64
	err := pr.DB.Model(&models.Pictures{}).Count(&count).Error
	return count, err
}

func (pr *PictureRepository) GetPicturesPaginated(lastSeenID int, limit int) ([]models.Pictures, error) {
	var pictures []models.Pictures
	err := pr.DB.Order("id").Limit(limit).Where("id > ?", lastSeenID).Find(&pictures).Error
	return pictures, err
}
