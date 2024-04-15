package models

type User struct {
	ID         int    `json:"id"`
	Permission int    `json:"permission"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
}
