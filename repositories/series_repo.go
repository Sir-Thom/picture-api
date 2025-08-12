package repositories

import (
	"Api-Picture/models"
	"gorm.io/gorm"
)

type SeriesRepository struct {
	DB *gorm.DB
}

func NewSeriesRepository(db *gorm.DB) *SeriesRepository {
	return &SeriesRepository{DB: db}
}

func (sr *SeriesRepository) CreateSeries(series *models.Series) error {
	return sr.DB.Create(series).Error
}

func (sr *SeriesRepository) GetAllSeries() ([]models.Series, error) {
	var series []models.Series
	err := sr.DB.Preload("Video").Find(&series).Error
	return series, err
}

func (sr *SeriesRepository) GetSeriesByID(id uint) (*models.Series, error) {
	var series models.Series
	//Preload("Video") is used to load the associated videos of the series
	err := sr.DB.Preload("Video").First(&series, id).Error
	return &series, err
}

func (sr *SeriesRepository) GetSeriesByName(name string) ([]models.Series, error) {
	var series []models.Series
	err := sr.DB.Preload("Video").Where("name ILIKE ?", "%"+name+"%").Find(&series).Error
	return series, err
}
