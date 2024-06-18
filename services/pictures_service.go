package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
)

type PictureService struct {
	Repository *repositories.PictureRepository
}

func NewPictureService(repo *repositories.PictureRepository) *PictureService {
	return &PictureService{Repository: repo}
}

func (ps *PictureService) GetAllPictures(limit int) ([]models.Pictures, error) {
	return ps.Repository.GetAll(limit)
}

func (ps *PictureService) GetPictureById(id string) (models.Pictures, error) {
	return ps.Repository.GetById(id)
}

func (ps *PictureService) CountPictures() (int64, error) {
	return ps.Repository.Count()
}

func (ps *PictureService) GetPicturesPaginated(lastSeenID int, limit int) ([]models.Pictures, error) {
	const batchSize = 100
	return ps.Repository.GetPicturesPaginated(lastSeenID, limit, batchSize)
}
