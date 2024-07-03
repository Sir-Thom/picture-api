package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
)

type VideoService struct {
	Repository *repositories.VideoRepository
}

func NewVideoService(repo *repositories.VideoRepository) *VideoService {
	return &VideoService{Repository: repo}
}

func (vs *VideoService) GetAllVideo() ([]models.Video, error) {
	return vs.Repository.GetAllVideo()
}

func (vs *VideoService) GetVideoByID(id uint) (*models.Video, error) {
	return vs.Repository.GetVideoByID(id)
}

func (vs *VideoService) GetVideoByName(name string) ([]models.Video, error) {
	return vs.Repository.GetVideoByName(name)
}
