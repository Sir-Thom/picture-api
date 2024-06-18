package repositories

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate models
	err = db.AutoMigrate(&models.Pictures{})
	if err != nil {
		return nil
	}

	// Seed some data for testing
	for i := 0; i < 1000; i++ {
		db.Create(&models.Pictures{ID: int(uint(i + 1)), Filename: "image" + strconv.Itoa(i+1) + ".jpg", Data: []byte("data")})
	}

	return db
}

func TestGetAll(t *testing.T) {
	db := setupDB()

	repo := repositories.NewPictureRepository(db)

	pictures, err := repo.GetAll(10)
	assert.NoError(t, err)
	assert.Len(t, pictures, 10) // Assuming there are 3 pictures seeded
}

func TestGetById(t *testing.T) {
	db := setupDB()

	repo := repositories.NewPictureRepository(db)

	picture, err := repo.GetById("2")
	assert.NoError(t, err)
	assert.NotNil(t, picture)
	assert.Equal(t, "image2.jpg", picture.Filename)
}

func TestCount(t *testing.T) {
	db := setupDB()

	repo := repositories.NewPictureRepository(db)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.EqualValues(t, 1000, count) // Assuming there are 3 pictures seeded
}

// Write tests for GetPicturesPaginated function
func TestGetPicturesPaginated(t *testing.T) {
	db := setupDB()

	repo := repositories.NewPictureRepository(db)

	// Test pagination with batch size less than available pictures
	pictures, err := repo.GetPicturesPaginated(0, 2, 2)
	assert.NoError(t, err)
	assert.Len(t, pictures, 2)

	// Test pagination with batch size more than available pictures
	pictures, err = repo.GetPicturesPaginated(0, 10, 10)
	assert.NoError(t, err)
	assert.Len(t, pictures, 10) // Assuming there are 3 pictures seeded
}
func BenchmarkGetAll(b *testing.B) {
	db := setupDB()

	repo := repositories.NewPictureRepository(db)

	// Seed some data for testing

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetAll(1000)
		assert.NoError(b, err)
	}
}
