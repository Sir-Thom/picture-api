package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"github.com/dgrijalva/jwt-go"
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

func (us *UserService) SignUp(email, password string) (error, error) {
	user := models.User{
		Email:    email,
		Password: password,
	}

	return us.Repo.SignUp(user), nil

}

func NewJWTService(secretKey string, tokenExpiration time.Duration) *JWTService {
	return &JWTService{SecretKey: secretKey, TokenExpiration: tokenExpiration}
}

func (j *JWTService) GenerateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(j.TokenExpiration).Unix()

	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

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
