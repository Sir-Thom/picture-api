package services

import (
	"Api-Picture/models"
	"Api-Picture/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	SecretKey       string
	TokenExpiration time.Duration
}

type UserService struct {
	Repo       *repositories.UserRepository
	JWTService *JWTService
}

func NewUserService(repo *repositories.UserRepository, jwtService *JWTService) *UserService {
	return &UserService{Repo: repo, JWTService: jwtService}
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

	token, err := us.JWTService.GenerateToken(user.ID)
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

	token, err := us.JWTService.GenerateToken(user.ID)
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
	expirationTime := time.Now().Add(j.TokenExpiration)
	claims["exp"] = expirationTime.Unix()

	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.SecretKey), nil
	})
}
