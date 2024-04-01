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
		// If no more pictures are found, break the loop
		if len(batch) == 0 {
			break
		}
		pictures = append(pictures, batch...)
		lastSeenID = batch[len(batch)-1].ID // Update lastSeenID for the next batch
	}

	// If there are more than zero but less than the limit pictures left, create a new page with the last picture
	if len(pictures) > 0 && len(pictures) < limit {
		lastPicture := pictures[len(pictures)-1]
		fmt.Println("Last Picture ID:", lastPicture.ID) // Log lastPicture.ID
		newPage, err := pr.GetPicturesPaginated(lastPicture.ID, limit, batchSize)
		if err != nil {
			return nil, err
		}
		fmt.Println("New Page Length:", len(newPage)) // Log length of newPage
		if len(newPage) > 0 {
			pictures = append(pictures, newPage...)
		} else {
			// Fetch the last image using the last seen ID and append it to the list of pictures
			lastImage, err := pr.GetById(fmt.Sprintf("%d", lastSeenID-limit))
			if err != nil {
				return nil, err
			}
			pictures = append(pictures, lastImage)
		}
	}

	// Cache the paginated pictures
	pr.Cache.Set(cacheKey, pictures, cache.DefaultExpiration)

	return pictures, nil
}
