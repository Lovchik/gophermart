package handlers

import (
	"fmt"
	"github.com/Lovchik/gophermart/internal/server/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (s *Service) CreateOrders(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	c.Request.Body.Close()
	orderNumber := string(bodyBytes)
	userID, err := s.JwtKeysPair.GetUserID(c.GetHeader("Authorization"))
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
	if !validateOrderNumber(orderNumber) {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	info, err := s.Feign.GetBonusInfo(orderNumber)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	err = s.Store.CreateOrder(orderNumber, userID, info)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)
}

func (s *Service) GetOrders(context *gin.Context) {
	userID, err := s.JwtKeysPair.GetUserID(context.GetHeader("Authorization"))
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
		log.Println(err)
		context.Status(http.StatusBadRequest)
		return
	}
	err = s.Validator.Struct(order)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	if !validateOrderNumber(order.Order) {
		context.Status(http.StatusUnprocessableEntity)
		return
	}
	userID, err := s.JwtKeysPair.GetUserID(context.GetHeader("Authorization"))
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
	userID, err := s.JwtKeysPair.GetUserID(context.GetHeader("Authorization"))
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
	userID, err := s.JwtKeysPair.GetUserID(context.GetHeader("Authorization"))
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
	var sum int
	var alt bool
	for i := len(orderNumber) - 1; i >= 0; i-- {
		r := orderNumber[i]

		if r < '0' || r > '9' {
			return false
		}
		n := int(r - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}
