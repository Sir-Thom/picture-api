package models

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}
