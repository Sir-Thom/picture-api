package repositories

import (
	"Api-Picture/models"

	"gorm.io/gorm"
)

type VideoRepository struct {
	DB *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{DB: db}
}

// GetVideoByID retrieves single video by ID
func (vr *VideoRepository) GetVideoByID(id uint) (*models.Video, error) {
	var video models.Video
	err := vr.DB.First(&video, id).Error
	return &video, err
}

// GetVideosPaginated retrieves paginated videos
func (vr *VideoRepository) GetVideosPaginated(offset, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	// Count total records
	if err := vr.DB.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err := vr.DB.Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetVideoByNamePaginated searches videos by name with pagination
func (vr *VideoRepository) GetVideoByNamePaginated(name string, offset, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	// Count total matches
	query := vr.DB.Model(&models.Video{}).Where("video_name ILIKE ?", "%"+name+"%")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err := query.Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetVideosBySeriesID retrieves videos by series ID
func (vr *VideoRepository) GetVideosBySeriesID(seriesID uint) ([]models.Video, error) {
	var videos []models.Video
	err := vr.DB.Where("series_id = ?", seriesID).
		Order("id ASC").
		Find(&videos).Error
	return videos, err
}

// GetVideoMetadataOnly retrieves lightweight video metadata
func (vr *VideoRepository) GetVideoMetadataOnly() ([]models.Video, error) {
	var videos []models.Video
	err := vr.DB.Select("id, video_name, series_id").
		Order("id DESC").
		Find(&videos).Error
	return videos, err
}

// GetVideosByNameExact finds videos by exact name match
func (vr *VideoRepository) GetVideosByNameExact(name string) ([]models.Video, error) {
	var videos []models.Video
	err := vr.DB.Where("video_name = ?", name).
		Find(&videos).Error
	return videos, err
}

// BulkGetVideosByIDs retrieves multiple videos by their IDs
func (vr *VideoRepository) BulkGetVideosByIDs(ids []uint) ([]models.Video, error) {
	var videos []models.Video
	err := vr.DB.Where("id IN ?", ids).
		Order("id DESC").
		Find(&videos).Error
	return videos, err
}
