package models

type User struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Username string `gorm:"not null"`
	Password string `gorm:"not null"`
}
