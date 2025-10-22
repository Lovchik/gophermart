package handlers

import (
	"fmt"
	"github.com/Lovchik/gophermart/internal/server/models"
	"github.com/Lovchik/gophermart/internal/server/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func (s *Service) CreateOrders(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	c.Request.Body.Close()
	orderNumber := string(bodyBytes)

	if !validateOrderNumber(orderNumber) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := utils.GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	ownerID, exists := s.Store.GetOrderOwner(orderNumber)
	if exists {
		if ownerID == userID {
			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusConflict)
		return
	}
	err = s.Store.CreateOrder(orderNumber, userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)
}

func (s *Service) GetOrders(context *gin.Context) {
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}
	orders, err := s.Store.GetOrders(userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		context.Status(http.StatusNoContent)
		return
	}

	context.JSON(http.StatusOK, orders)
}

func (s *Service) CreateWithdraw(context *gin.Context) {
	var order models.CreateWithdrawalOrder
	err := context.ShouldBind(&order)
	if err != nil {
		fmt.Println(err)
		context.Status(http.StatusBadRequest)
		return
	}
	err = utils.Validate().Struct(order)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	if validateOrderNumber(order.Order) {
		context.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}
	balance, err := s.Store.GetActualBalance(userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	if balance-order.Sum < 0 {
		context.Status(http.StatusPaymentRequired)
		return
	}
	err = s.Store.CreateWithdrawalOrder(order, userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}

func (s *Service) GetBalance(context *gin.Context) {
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}
	withdraw, err := s.Store.GetWithdraw(userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}
	bonuses, err := s.Store.GetBonuses(userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"current": bonuses - withdraw, "withdrawn": withdraw})
}

func (s *Service) GetWithdrawals(context *gin.Context) {
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.Status(http.StatusUnauthorized)
		return
	}
	orders, err := s.Store.GetWithdrawalOrders(userID)
	if err != nil {
		fmt.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		context.Status(http.StatusNoContent)
		return
	}

	context.JSON(http.StatusOK, orders)
}

func validateOrderNumber(orderNumber string) bool {
	var result int64
	for i, char := range orderNumber {
		parseInt, err := strconv.ParseInt(string(char), 10, 64)
		if err != nil {
			return false
		}
		if i+1%2 == 0 {
			parseInt = parseInt * 2
		}
		if parseInt > 9 {
			parseInt = parseInt - 9
		}
		result += parseInt
	}

	return result%10 == 0
}
