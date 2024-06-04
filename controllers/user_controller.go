package controllers

import (
	"Api-Picture/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// SignUp godoc
// @Summary Sign up
// @Description Sign up
// @Tags users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param username formData string true "Username"
// @Success 200 {string} string "ok"
// @Router /signup/register [post]
func (uc *UserController) SignUp(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	username := ctx.PostForm("username")

	err, token := uc.service.SignUp(email, password, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "token": token})
}

// SignIn godoc
// @Summary Sign in
// @Description Sign in
// @Tags users
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {string} string "ok"
// @Router /signin [post]
func (uc *UserController) SignIn(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	err, token := uc.service.SignIn(email, password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "token": token})
}
