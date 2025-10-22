package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gofermart/internal/server/models"
	"gofermart/internal/server/storage"
	"gofermart/internal/server/utils"
	"net/http"
)

type Service struct {
	WebServer *gin.Engine
	Store     storage.Storage
}

func (s *Service) Refresh(c *gin.Context) {
	var response models.Response
	header := c.GetHeader("Refresh")
	if !utils.IsValidToken(header, "refresh") {
		log.Info("Неверный токен")
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Неверный токен"))
		return
	}
	id, err := utils.GetUserID(header)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	tokens, err := utils.GenerateJWT(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Authorization", tokens.AccessToken)
	c.Header("Refresh", tokens.RefreshToken)
	c.JSON(http.StatusOK, nil)
}

func (s *Service) Login(c *gin.Context) {
	var response models.Response
	var credentials models.LoginRequest
	err := c.ShouldBind(&credentials)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Неверные параметры тела запроса"))
		return
	}
	err = utils.Validate().Struct(credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Должны быть заполнены логин и пароль"))
		return
	}
	user, err := s.Store.GetUserByCreds(credentials)
	if err != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Неверная пара логин/пароль"))
		return
	}
	tokens, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	c.Header("Authorization", tokens.AccessToken)
	c.Header("Refresh", tokens.RefreshToken)
	c.JSON(http.StatusOK, nil)
}

func (s *Service) RegisterUser(c *gin.Context) {
	var response models.Response
	var request models.LoginRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	exists := s.Store.IsUserExists(request)
	if exists {
		c.JSON(http.StatusConflict, nil)
		return
	}

	user, err := s.Store.CreateUser(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	tokens, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.Header("Authorization", tokens.AccessToken)
	c.Header("Refresh", tokens.RefreshToken)
	c.JSON(http.StatusOK, response.NewWithMessage(request, ""))
}
