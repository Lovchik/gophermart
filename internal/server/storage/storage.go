package storage

import (
	"context"
	"github.com/Lovchik/gophermart/internal/server/feign"
	"github.com/Lovchik/gophermart/internal/server/models"
)

type Storage interface {
	HealthCheck() error
	Close(ctx context.Context) error
	IsUserExists(request models.LoginRequest) bool
	CreateUser(request models.LoginRequest) (models.User, error)
	GetOrderOwner(orderNumber string) (int64, bool)
	CreateOrder(orderNumber string, userID int64, info feign.BonusInfo) error
	GetOrders(userID int64) ([]models.Order, error)
	GetWithdrawalOrders(userID int64) ([]models.WithdrawalOrders, error)
	GetUserByCreds(user models.LoginRequest) (models.User, error)
	CreateWithdrawalOrder(order models.CreateWithdrawalOrder, userID int64) error
	GetBonuses(userID int64) (float64, error)
	GetWithdraw(userID int64) (float64, error)
	GetActualBalance(userID int64) (float64, error)
}
