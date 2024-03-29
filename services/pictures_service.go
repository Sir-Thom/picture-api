package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
)

type PictureService struct {
	Repo *repositories.PictureRepository
}

func NewPictureService(repo *repositories.PictureRepository) *PictureService {
	return &PictureService{Repo: repo}
}

func (ps *PictureService) GetAllPictures(limit int) ([]models.Pictures, error) {
	return ps.Repo.GetAll(limit)
}

func (ps *PictureService) GetPictureById(id string) (models.Pictures, error) {
	return ps.Repo.GetById(id)
}

func (ps *PictureService) CountPictures() (int64, error) {
	return ps.Repo.Count()
}

func (ps *PictureService) GetPicturesAfterID(lastSeenID int, limit int) ([]models.Pictures, error) {

	return ps.Repo.GetPicturesAfterID(lastSeenID, limit)
}

func (ps *PictureService) GetPicturesPaginated(lastSeenID int, limit int) ([]models.Pictures, error) {
	return ps.Repo.GetPicturesPaginated(lastSeenID, limit)
}
