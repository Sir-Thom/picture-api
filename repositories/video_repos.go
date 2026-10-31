package repositories

import (
	"Api-Picture/models"
	"time"

	lru "github.com/hashicorp/golang-lru/v2/expirable"
	"gorm.io/gorm"
)

type VideoRepository struct {
	DB *gorm.DB
	// paths caches id -> video_path. A <video> element issues many Range
	// requests while playing/seeking; this keeps them off the database.
	paths *lru.LRU[uint, string]
	// counts caches COUNT(*) results used for pagination metadata.
	counts *lru.LRU[string, int64]
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{
		DB:     db,
		paths:  lru.NewLRU[uint, string](2048, nil, 10*time.Minute),
		counts: lru.NewLRU[string, int64](512, nil, time.Minute),
	}
}

func (vr *VideoRepository) GetVideoPathByID(id uint) (string, error) {
	if p, ok := vr.paths.Get(id); ok {
		return p, nil
	}
	var video models.Video
	if err := vr.DB.Select("id, video_path").First(&video, id).Error; err != nil {
		return "", err
	}
	vr.paths.Add(id, video.VideoPath)
	return video.VideoPath, nil
}

func (vr *VideoRepository) cachedCount(key string, q *gorm.DB) (int64, error) {
	if n, ok := vr.counts.Get(key); ok {
		return n, nil
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return 0, err
	}
	vr.counts.Add(key, n)
	return n, nil
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
	//var total int64

	// Count total records (cached: COUNT(*) is a full scan in PostgreSQL)
	total, err := vr.cachedCount("all", vr.DB.Model(&models.Video{}))
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err = vr.DB.Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetVideoByNamePaginated searches videos by name with pagination
func (vr *VideoRepository) GetVideoByNamePaginated(name string, offset, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	// Count total matches. ILIKE '%x%' is served by the pg_trgm GIN index
	// created in models.EnsureSearchIndexes.
	pattern := "%" + name + "%"
	total, err := vr.cachedCount("name:"+name,
		vr.DB.Model(&models.Video{}).Where("video_name ILIKE ?", pattern))
	if err != nil {
		return nil, 0, err
	}
	// Fetch paginated results (fresh query: reusing a Count()ed *gorm.DB
	// carries its state over and is a common source of subtle bugs)
	err = vr.DB.Where("video_name ILIKE ?", pattern).
		Order("id DESC").
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
