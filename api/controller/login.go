package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/companies/api/dto"
	"github.com/kumareswaramoorthi/companies/api/errors"
	service "github.com/kumareswaramoorthi/companies/api/service"
)

type LoginController interface {
	Login(ctx *gin.Context)
}

type loginController struct {
	loginService service.LoginService
	jWtService   service.JWTService
}

func NewLoginController(loginService service.LoginService,
	jWtService service.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(c *gin.Context) {
	var credential dto.LoginCredentials
	err := c.ShouldBind(&credential)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}
	isUserAuthenticated := controller.loginService.LoginUser(credential.Email, credential.Password)
	if !isUserAuthenticated {
		c.AbortWithStatusJSON(http.StatusUnauthorized, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": controller.jWtService.GenerateToken(credential.Email, true),
	})
}
