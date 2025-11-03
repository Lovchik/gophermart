package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(router *gin.RouterGroup, s *Service) {
	router.POST("/refresh", s.Refresh)
	router.POST("/login", s.Login)
	router.POST("/register", s.RegisterUser)
	router.POST("/orders", s.AuthMiddleware(), s.CreateOrders)
	router.GET("/orders", s.AuthMiddleware(), s.GetOrders)
	router.GET("/balance", s.AuthMiddleware(), s.GetBalance)
	router.POST("/balance/withdraw", s.AuthMiddleware(), s.CreateWithdraw)
	router.GET("/withdrawals", s.AuthMiddleware(), s.GetWithdrawals)
}

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token, _ = c.Cookie("Authorization")
		}
		if token == "" && !s.JwtKeysPair.IsValidToken(token, "access") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := s.JwtKeysPair.GetUserID(token)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
