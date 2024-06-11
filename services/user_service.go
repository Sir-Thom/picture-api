package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
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

func (us *UserService) SignUp(email, password, username string) (error, string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err, ""
	}

	user := models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	err = us.Repo.SignUp(user)
	if err != nil {
		return err, ""
	}

	jwtService := NewJWTService(os.Getenv("SECRET_KEY"), time.Hour)
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		return err, ""
	}

	return nil, token
}

func (us *UserService) SignIn(email, password string) (error, string) {
	user, err := us.Repo.Login(email)
	if err != nil {
		return err, ""
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err, ""
	}

	jwtService := NewJWTService(os.Getenv("SECRET_KEY"), time.Hour)
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		return err, ""
	}

	return nil, token
}

func NewJWTService(secretKey string, tokenExpiration time.Duration) *JWTService {
	return &JWTService{SecretKey: secretKey, TokenExpiration: tokenExpiration}
}

func (j *JWTService) GenerateToken(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	// Set expiration time to 90 days
	expirationTime := time.Now().Add(90 * 24 * time.Hour)
	claims["exp"] = expirationTime.Unix()

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
