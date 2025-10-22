package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gofermart/internal/server/models"
	"gofermart/internal/server/utils"
	"io"
	"net/http"
	"strconv"
)

func (s *Service) CreateOrders(c *gin.Context) {
	var response models.Response
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Неверный формат запроса"))
	}
	c.Request.Body.Close()
	orderNumber := string(bodyBytes)

	if !validateOrderNumber(orderNumber) {
		c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse("Неверный формат запроса"))
		return
	}
	userID, err := utils.GetUserID(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	ownerID, exists := s.Store.GetOrderOwner(orderNumber)
	if exists {
		if ownerID == userID {
			c.JSON(http.StatusOK, response.NewWithMessage(nil, "номер заказа уже был загружен этим пользователем"))
			return
		}
		c.JSON(http.StatusConflict, response.NewWithMessage(nil, "номер заказа уже был загружен другим пользователем"))
		return
	}
	err = s.Store.CreateOrder(orderNumber, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	c.JSON(http.StatusAccepted, response.NewWithMessage(nil, "новый номер заказа принят в обработку"))
}

func (s *Service) GetOrders(context *gin.Context) {
	var response models.Response
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse("Не авторизированный пользователь"))
		return
	}
	orders, err := s.Store.GetOrders(userID)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	if len(orders) == 0 {
		context.JSON(http.StatusNoContent, response.NewWithMessage(orders, "Успешно получено"))
		return
	}

	context.JSON(http.StatusOK, response.NewWithMessage(orders, "Успешно получено"))
}

func (s *Service) CreateWithdraw(context *gin.Context) {
	var response models.Response
	var order models.CreateWithdrawalOrder
	err := context.ShouldBind(&order)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, response.ErrorResponse("Неверные параметры тела запроса"))
		return
	}
	err = utils.Validate().Struct(order)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.ErrorResponse("Поля обязательны"))
		return
	}
	if validateOrderNumber(order.Order) {
		context.JSON(http.StatusUnprocessableEntity, response.ErrorResponse("Неверный формат запроса"))
		return
	}
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse("Пользователь не авторизован"))
		return
	}
	balance, err := s.Store.GetActualBalance(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	if balance-order.Sum < 0 {
		context.JSON(http.StatusPaymentRequired, response.ErrorResponse("на счету недостаточно средств"))
		return
	}
	err = s.Store.CreateWithdrawalOrder(order, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	context.JSON(http.StatusOK, response.ErrorResponse("успешная обработка запроса"))
}

func (s *Service) GetBalance(context *gin.Context) {
	var response models.Response
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse("Пользователь не авторизован"))
		return
	}
	withdraw, err := s.Store.GetWithdraw(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Серверная ошибка"))
		return
	}
	bonuses, err := s.Store.GetBonuses(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Серверная ошибка"))
		return
	}

	context.JSON(http.StatusOK, gin.H{"current": bonuses - withdraw, "withdrawn": withdraw})
}

func (s *Service) GetWithdrawals(context *gin.Context) {
	var response models.Response
	userID, err := utils.GetUserID(context.GetHeader("Authorization"))
	if err != nil {
		context.JSON(http.StatusUnauthorized, response.ErrorResponse("Не авторизированный пользователь"))
		return
	}
	orders, err := s.Store.GetWithdrawalOrders(userID)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, response.ErrorResponse("Внутренняя ошибка сервера."))
		return
	}
	if len(orders) == 0 {
		context.JSON(http.StatusNoContent, response.NewWithMessage(orders, "Успешно получено"))
		return
	}

	context.JSON(http.StatusOK, response.NewWithMessage(orders, "Успешно получено"))
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
