package handlers

import (
	"fmt"
	"github.com/Lovchik/gophermart/internal/server/models"
	"github.com/Lovchik/gophermart/internal/server/storage"
	"github.com/Lovchik/gophermart/internal/server/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Service struct {
	WebServer *gin.Engine
	Store     storage.Storage
}

func (s *Service) Refresh(c *gin.Context) {
	header := c.GetHeader("Refresh")
	if !utils.IsValidToken(header, "refresh") {
		log.Info("Неверный токен")
		c.Status(http.StatusUnauthorized)
		return
	}
	id, err := utils.GetUserID(header)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	tokens, err := utils.GenerateJWT(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Header("Authorization", tokens.AccessToken)
	c.Header("Refresh", tokens.RefreshToken)
	c.JSON(http.StatusOK, nil)
}

func (s *Service) Login(c *gin.Context) {
	var credentials models.LoginRequest
	err := c.ShouldBind(&credentials)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}
	err = utils.Validate().Struct(credentials)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	user, err := s.Store.GetUserByCreds(credentials)
	if err != nil || user.ID == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}
	tokens, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Authorization", tokens.AccessToken)
	c.Header("Refresh", tokens.RefreshToken)
	c.Status(http.StatusOK)
}

func (s *Service) RegisterUser(c *gin.Context) {
	var request models.LoginRequest
	err := c.ShouldBind(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	exists := s.Store.IsUserExists(request)
	if exists {
		c.Status(http.StatusConflict)
		return
	}

	_, err = s.Store.CreateUser(request)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}
