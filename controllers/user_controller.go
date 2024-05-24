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
//
// @Summary		Sign up
// @Description	Sign up
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			email body string true "Email"
// @Param			password body string true "Password"
// @Param			username body string true "Username"
// @Success		200 {string} string	"ok"
// @Router			/signup/register [post]
func (uc *UserController) SignUp(ctx *gin.Context) {

	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	username := ctx.PostForm("username")
	err, _ := uc.service.SignUp(email, password, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})

}
