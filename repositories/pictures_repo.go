package repositories

import (
	"Api-Picture/models"
	"fmt"
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
	cacheKey := fmt.Sprintf("all_pictures_%d", limit)
	if pictures, found := pr.Cache.Get(cacheKey); found {
		return pictures.([]models.Pictures), nil
	}

	var pictures []models.Pictures
	err := pr.DB.Select("id, filename, data").Limit(limit).Find(&pictures).Error
	if err == nil {
		pr.Cache.Set(cacheKey, pictures, cache.DefaultExpiration)
	}
	return pictures, err
}

func (pr *PictureRepository) GetById(id string) (models.Pictures, error) {
	cacheKey := fmt.Sprintf("picture_%s", id)
	if picture, found := pr.Cache.Get(cacheKey); found {
		return picture.(models.Pictures), nil
	}

	var picture models.Pictures
	err := pr.DB.Where("id = ?", id).Select("id, filename, data").First(&picture).Error
	if err == nil {
		pr.Cache.Set(cacheKey, picture, cache.DefaultExpiration)
	}
	return picture, err
}

func (pr *PictureRepository) Count() (int64, error) {
	const cacheKey = "pictures_count"
	if count, found := pr.Cache.Get(cacheKey); found {
		return count.(int64), nil
	}

	var count int64
	err := pr.DB.Model(&models.Pictures{}).Count(&count).Error
	if err == nil {
		pr.Cache.Set(cacheKey, count, cache.DefaultExpiration)
	}
	return count, err
}

func (pr *PictureRepository) GetPicturesPaginated(lastSeenID int, limit int, batchSize int) ([]models.Pictures, error) {
	cacheKey := fmt.Sprintf("paginated_pictures_%d_%d_%d", lastSeenID, limit, batchSize)

	if pictures, found := pr.Cache.Get(cacheKey); found {
		return pictures.([]models.Pictures), nil
	}

	var pictures []models.Pictures
	err := pr.DB.Order("id").Where("id > ?", lastSeenID).Limit(limit).Find(&pictures).Error
	if err != nil {
		return nil, err
	}

	pr.Cache.Set(cacheKey, pictures, cache.DefaultExpiration)

	return pictures, nil
}
