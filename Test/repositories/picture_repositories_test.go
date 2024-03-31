package repositories_test

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"time"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Order(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Limit(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return m
}

func (m *MockDB) Find(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Select(query interface{}) *gorm.DB {
	args := m.Called(query)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func GetAllPicturesReturnsFromCache(t *testing.T) {
	mockDB := new(MockDB)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	pictures := []models.Pictures{{ID: 1}, {ID: 2}}
	DefaultExpiration := 5 * time.Minute

	cache.Set("all_pictures", pictures, DefaultExpiration)

	repo := repositories.NewPictureRepository(mockDB)
	repo.Cache = cache

	result, _ := repo.GetAll(2)

	assert.Equal(t, pictures, result)
}

func GetByIdReturnsPictureFromDB(t *testing.T) {
	mockDB := new(MockDB)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	picture := models.Pictures{ID: 1}

	mockDB.On("Where", "id=?", []interface{}{"1"}).Return(mockDB)
	mockDB.On("Select", "id, filename, data").Return(mockDB)
	mockDB.On("First", &picture).Return(mockDB)

	repo := repositories.NewPictureRepository(mockDB)
	repo.Cache = cache

	result, _ := repo.GetById("1")

	assert.Equal(t, picture, result)
}

func CountReturnsCountFromDB(t *testing.T) {
	mockDB := new(MockDB)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	count := int64(2)

	mockDB.On("Model", &models.Pictures{}).Return(mockDB)
	mockDB.On("Count", &count).Return(mockDB)

	repo := repositories.NewPictureRepository(mockDB)
	repo.Cache = cache

	result, _ := repo.Count()

	assert.Equal(t, count, result)
}

func GetPicturesPaginatedReturnsFromDB(t *testing.T) {
	mockDB := new(MockDB)
	cache := cache.New(5*time.Minute, 10*time.Minute)
	pictures := []models.Pictures{{ID: 1}, {ID: 2}}

	mockDB.On("Order", "id").Return(mockDB)
	mockDB.On("Limit", 2).Return(mockDB)
	mockDB.On("Where", "id > ?", []interface{}{0}).Return(mockDB)
	mockDB.On("Find", &pictures).Return(mockDB)

	repo := repositories.NewPictureRepository(mockDB)
	repo.Cache = cache

	result, _ := repo.GetPicturesPaginated(0, 2)

	assert.Equal(t, pictures, result)
}
