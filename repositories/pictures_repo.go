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

func (pr *PictureRepository) GetPicturesPaginated(lastSeenID int, limit int, batchSize int) ([]models.Pictures, error) {
	// Generate a unique cache key based on lastSeenID, limit, and batchSize
	cacheKey := fmt.Sprintf("paginated_pictures_%d_%d_%d", lastSeenID, limit, batchSize)

	// Check if the paginated pictures are cached
	if pictures, found := pr.Cache.Get(cacheKey); found {
		return pictures.([]models.Pictures), nil
	}

	var pictures []models.Pictures

	// Fetch multiple pages (or a batch) of images in a single query
	for offset := 0; len(pictures) < batchSize; offset += limit {
		var batch []models.Pictures
		err := pr.DB.Order("id").Limit(limit).Offset(offset).Where("id > ?", lastSeenID).Find(&batch).Error
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			// No more pictures available
			break
		}
		pictures = append(pictures, batch...)
		lastSeenID = batch[len(batch)-1].ID // Update lastSeenID for the next batch
	}

	// Cache the paginated pictures
	pr.Cache.Set(cacheKey, pictures, cache.DefaultExpiration)

	return pictures, nil
}
