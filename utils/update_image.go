package utils

import (
	"Api-Picture/models"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

func UpdateDb() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovering from panic: %v", r)
			}
		}()

		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		path := os.Getenv("FOLDER_PATH")

		db, err := models.Database()
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		err = db.AutoMigrate(&models.Pictures{})
		if err != nil {
			log.Fatalf("Error migrating database: %v", err)
		}

		var folderFilenames, dbFilenames []string
		done := make(chan struct{})
		go func() {
			defer close(done)
			getFilesInFolder(path, &folderFilenames)
		}()
		go func() {
			defer close(done)
			getDatabasePictures(db, &dbFilenames)
		}()

		for i := 0; i < 2; i++ {
			<-done
		}

		files, err := os.ReadDir(path)
		if err != nil {
			log.Fatalf("Error reading directory: %v", err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".webp") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".gif") {
				if contains(dbFilenames, file.Name()) {
					go updatePicture(db, path, file)
				} else {
					go insertPicture(db, path, file)
				}
			}
		}

		go deletePicturesNotInFolder(db, dbFilenames, folderFilenames)

		log.Println("Database update started")
	}()
}

func getFilesInFolder(path string, folderFilenames *[]string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".webp") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".gif") {
			*folderFilenames = append(*folderFilenames, file.Name())
		}
	}
}

func getDatabasePictures(db *gorm.DB, dbFilenames *[]string) {
	var pictures []models.Pictures
	db.Find(&pictures)

	for _, pic := range pictures {
		log.Println(pic.Filename)
		*dbFilenames = append(*dbFilenames, pic.Filename)
	}
}

func updatePicture(db *gorm.DB, path string, file os.DirEntry) {
	data, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		log.Printf("Error reading file %s: %v", file.Name(), err)
		return
	}

	var picture models.Pictures
	db.Where("filename = ?", file.Name()).First(&picture)
	picture.Data = data
	picture.AddedDate = time.Now()
	db.Save(&picture)
}

func insertPicture(db *gorm.DB, path string, file os.DirEntry) {
	data, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		log.Printf("Error reading file %s: %v", file.Name(), err)
		return
	}

	picture := models.Pictures{
		Filename:  file.Name(),
		Data:      data,
		AddedDate: time.Now(),
	}
	db.Create(&picture)
}

func deletePicturesNotInFolder(db *gorm.DB, dbFilenames []string, folderFilenames []string) {
	for _, dbFilename := range dbFilenames {
		if !contains(folderFilenames, dbFilename) {
			db.Where("filename = ?", dbFilename).Delete(&models.Pictures{})
		}
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
