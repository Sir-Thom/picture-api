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

func (vs *VideoService) GetVideosPaginated(offset, limit int) ([]models.Video, int64, error) {
	return vs.Repository.GetVideosPaginated(offset, limit)
}

func (vs *VideoService) GetVideoByNamePaginated(name string, offset, limit int) ([]models.Video, int64, error) {
	return vs.Repository.GetVideoByNamePaginated(name, offset, limit)
}

func (vs *VideoService) GetVideoByID(id uint) (*models.Video, error) {
	return vs.Repository.GetVideoByID(id)
}

// Additional optimized service methods
func (vs *VideoService) GetVideosBySeriesID(seriesID uint) ([]models.Video, error) {
	return vs.Repository.GetVideosBySeriesID(seriesID)
}

func (vs *VideoService) GetVideoMetadataOnly() ([]models.Video, error) {
	return vs.Repository.GetVideoMetadataOnly()
}

func (vs *VideoService) GetVideosByNameExact(name string) ([]models.Video, error) {
	return vs.Repository.GetVideosByNameExact(name)
}

func (vs *VideoService) BulkGetVideosByIDs(ids []uint) ([]models.Video, error) {
	return vs.Repository.BulkGetVideosByIDs(ids)
}
