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

func (vr *VideoRepository) GetVideoByID(id uint) (*models.Video, error) {
	var video models.Video
	err := vr.DB.First(&video, id).Error
	return &video, err
}

func (vr *VideoRepository) GetVideoByName(name string) ([]models.Video, error) {
	var videos []models.Video
	err := vr.DB.Where("video_name ILIKE ?", "%"+name+"%").Find(&videos).Error
	return videos, err
}

func (vr *VideoRepository) GetAllVideo() ([]models.Video, error) {
	var video []models.Video
	err := vr.DB.Select("id, video_name, series_id,video_path").Find(&video).Error
	return video, err
}
