package models

type Video struct {
	ID        uint   `gorm:"primaryKey"`
	SeriesID  uint   `gorm:"not null"`
	VideoName string `gorm:"size:255;not null"`
	VideoPath string `gorm:"size:255;not null"`
}
