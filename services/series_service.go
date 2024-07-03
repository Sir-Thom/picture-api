package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
)

type SeriesService struct {
	Repository *repositories.SeriesRepository
}

func NewSeriesService(repo *repositories.SeriesRepository) *SeriesService {
	return &SeriesService{Repository: repo}
}

func (ss *SeriesService) GetAllSeries() ([]models.Series, error) {
	return ss.Repository.GetAllSeries()
}

func (ss *SeriesService) GetSeriesByID(id uint) (*models.Series, error) {
	return ss.Repository.GetSeriesByID(id)
}

func (ss *SeriesService) GetSeriesByName(name string) ([]models.Series, error) {
	return ss.Repository.GetSeriesByName(name)
}
