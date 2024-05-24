package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type JWTService struct {
	SecretKey       string
	TokenExpiration time.Duration
}
type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (us *UserService) SignUp(email, password, username string) (error, error) {
	user := models.User{
		ID:       int(time.Now().Unix()),
		Email:    email,
		Username: username,
		Password: password,
	}
	w := NewJWTService("secret", time.Hour)
	log.Println(w)
	err := us.Repo.SignUp(user)
	if err != nil {
		return err, nil
	}
	token, err := w.GenerateToken(user.ID)
	if err != nil {
		return err, nil

	}
	log.Println(token)
	return nil, nil

}

func NewJWTService(secretKey string, tokenExpiration time.Duration) *JWTService {
	return &JWTService{SecretKey: secretKey, TokenExpiration: tokenExpiration}
}

func (j *JWTService) GenerateToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(j.TokenExpiration).Unix()

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	log.Println(tokenString)
	return tokenString, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil

}
