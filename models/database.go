package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

// EnsureSearchIndexes creates trigram indexes so the `ILIKE '%name%'`
// searches on videos and series can use an index instead of scanning every
// row (a plain B-tree index cannot serve a leading wildcard). It is
// best-effort: if the database user may not create the pg_trgm extension,
// a warning is logged and the API keeps working without the index.
func EnsureSearchIndexes(db *gorm.DB) {
	if db.Dialector.Name() != "postgres" {
		return
	}
	stmts := []string{
		"CREATE EXTENSION IF NOT EXISTS pg_trgm",
		"CREATE INDEX IF NOT EXISTS idx_videos_video_name_trgm ON videos USING gin (video_name gin_trgm_ops)",
		"CREATE INDEX IF NOT EXISTS idx_series_name_trgm ON series USING gin (name gin_trgm_ops)",
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("search index setup skipped (%s): %v", stmt, err)
			return
		}
	}
}
