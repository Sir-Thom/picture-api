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
// @Success		200 {string} string	"ok"
// @Router			/signup/register [post]
func (uc *UserController) SignUp(ctx *gin.Context) {
	signinUser, err := uc.service.SignUp("email", "password")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, signinUser)

}
