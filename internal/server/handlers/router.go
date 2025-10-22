package handlers

import (
	"fmt"
	"github.com/Lovchik/gophermart/internal/server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(router *gin.RouterGroup, s *Service) {
	router.POST("/refresh", s.Refresh)
	router.POST("/login", s.Login)
	router.POST("/register", s.RegisterUser)
	router.POST("/orders", AuthMiddleware(), s.CreateOrders)
	router.GET("/orders", AuthMiddleware(), s.GetOrders)
	router.GET("/balance", AuthMiddleware(), s.GetBalance)
	router.POST("/balance/withdraw", AuthMiddleware(), s.CreateWithdraw)
	router.GET("/withdrawals", AuthMiddleware(), s.GetWithdrawals)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		fmt.Println("token: ", token)
		if token == "" {
			token, _ = c.Cookie("Authorization")
		}
		if token == "" && !utils.IsValidToken(token, "access") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := utils.GetUserID(token)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}
