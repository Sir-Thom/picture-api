package repositories

import (
	"Api-Picture/models"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
	"gorm.io/gorm"
)

type PictureRepository struct {
	DB    *gorm.DB
	Cache *lru.LRU[string, any]
}

func NewPictureRepository(db *gorm.DB) *PictureRepository {
	// Holds up to 1000 items; items expire after 5 minutes
	cache := lru.NewLRU[string, any](1000, nil, 5*time.Minute)
	return &PictureRepository{
		DB:    db,
		Cache: cache,
	}
}

func (pr *PictureRepository) GetAll(limit int) ([]models.Pictures, error) {
	cacheKey := fmt.Sprintf("all_pictures_%d", limit)
	if pictures, found := pr.Cache.Get(cacheKey); found {
		return pictures.([]models.Pictures), nil
	}

	var pictures []models.Pictures
	err := pr.DB.Select("id, filename, path").Limit(limit).Find(&pictures).Error
	if err == nil {
		pr.Cache.Add(cacheKey, pictures)
	}
	return pictures, err
}

func (pr *PictureRepository) GetById(id string) (models.Pictures, error) {
	cacheKey := fmt.Sprintf("picture_%s", id)
	if picture, found := pr.Cache.Get(cacheKey); found {
		return picture.(models.Pictures), nil
	}

	var picture models.Pictures
	err := pr.DB.Where("id = ?", id).Select("id, filename, path").First(&picture).Error
	if err == nil {
		pr.Cache.Add(cacheKey, picture)
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
		pr.Cache.Add(cacheKey, count)
	}
	return count, err
}

func (pr *PictureRepository) GetPicturesPaginated(lastSeenID int, limit int, batchSize int) ([]models.Pictures, error) {
	cacheKey := fmt.Sprintf("paginated_pictures_%d_%d", lastSeenID, limit)

	if val, found := pr.Cache.Get(cacheKey); found {
		return val.([]models.Pictures), nil
	}

	var pictures []models.Pictures
	err := pr.DB.Order("id").Where("id > ?", lastSeenID).Limit(limit).Find(&pictures).Error
	if err != nil {
		return nil, err
	}

	pr.Cache.Add(cacheKey, pictures)
	return pictures, nil
}
