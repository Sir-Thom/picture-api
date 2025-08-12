package models

type Series struct {
	ID         uint    `gorm:"primaryKey"`
	Name       string  `gorm:"size:255;not null"`
	FolderPath string  `gorm:"size:255;not null"`
	Video      []Video `gorm:"foreignKey:SeriesID"`
}
