package storage

import (
	"context"
	"fmt"
	"github.com/Lovchik/gophermart/internal/server/feign"
	"github.com/Lovchik/gophermart/internal/server/models"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

type PostgresStorage struct {
	Conn *pgxpool.Pool
}

func (p PostgresStorage) GetWithdrawalOrders(userID int64) ([]models.WithdrawalOrders, error) {
	var orders []models.WithdrawalOrders
	query := `SELECT order_number, accrual,processed_at FROM withdrawal_orders WHERE user_id = $1 ORDER BY processed_at DESC`

	rows, err := p.Conn.Query(context.Background(), query, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.WithdrawalOrders
		if err := rows.Scan(&o.Order, &o.Accurual, &o.ProcessedAt); err != nil {
			fmt.Println(err)
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orders, nil
}

func (p PostgresStorage) GetActualBalance(userID int64) (float64, error) {
	withdraw, err := p.GetWithdraw(userID)
	if err != nil {
		return 0, err
	}
	bonuses, err := p.GetBonuses(userID)
	if err != nil {
		return 0, err
	}
	return bonuses - withdraw, nil
}

func (p PostgresStorage) GetBonuses(userID int64) (float64, error) {
	var bonusBalance float64
	bonusQuery := `SELECT SUM(orders.accrual) FROM orders WHERE user_id = $1`
	err := p.Conn.QueryRow(context.Background(), bonusQuery, userID).Scan(&bonusBalance)
	if err != nil {
		return 0, err
	}
	return bonusBalance, nil

}
func (p PostgresStorage) GetWithdraw(userID int64) (float64, error) {
	var withdrawBalance float64
	withdrawQuery := `SELECT SUM(withdrawal_orders.accrual) FROM withdrawal_orders WHERE user_id = $1`
	err := p.Conn.QueryRow(context.Background(), withdrawQuery, userID).Scan(&withdrawBalance)
	if err != nil {
		return 0, err
	}
	return withdrawBalance, nil
}

func (p PostgresStorage) CreateWithdrawalOrder(order models.CreateWithdrawalOrder, userID int64) error {
	var result int64
	query := `INSERT INTO withdrawal_orders (user_id, order_number,accrual) VALUES ($1, $2,$3)`
	err := p.Conn.QueryRow(context.Background(), query, userID, order.Order, order.Sum).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresStorage) GetUserByCreds(login models.LoginRequest) (models.User, error) {
	var user models.User
	query := `SELECT id,login FROM users WHERE login = $1 AND password = $2 LIMIT 1`
	err := p.Conn.QueryRow(context.Background(), query, login.Login, login.Password).Scan(&user.ID, &user.Login)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}
	return user, nil
}

func (p PostgresStorage) GetOrders(userID int64) ([]models.Order, error) {
	var orders []models.Order
	query := `SELECT number, status,accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`

	rows, err := p.Conn.Query(context.Background(), query, userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.Number, &o.Status, &o.Accurual, &o.UploadedAt); err != nil {
			fmt.Println(err)
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return orders, nil
}

func (p PostgresStorage) GetOrderOwner(orderNumber string) (int64, bool) {
	var userID int64
	query := `SELECT user_id FROM orders WHERE number = $1`
	err := p.Conn.QueryRow(context.Background(), query, orderNumber).Scan(&userID)
	if err != nil {
		return 0, false
	}
	return userID, true
}

func (p PostgresStorage) CreateOrder(orderNumber string, userID int64) error {
	var result int64
	var status string
	var accrual *float64
	info, err := feign.GetBonusInfo(orderNumber)
	if err != nil {
		return err
	}
	if info.Status == nil {
		status = "NEW"
	} else {
		status = *info.Status
	}
	if info.Accrual != nil {
		accrual = info.Accrual
	}

	query := `INSERT INTO orders (user_id, number,status,accrual) VALUES ($1,$2,$3,$4)`
	err = p.Conn.QueryRow(context.Background(), query, userID, orderNumber, status, accrual).Scan(&result)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresStorage) IsUserExists(request models.LoginRequest) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 AND password = $2)`
	err := p.Conn.QueryRow(context.Background(), query, request.Login, request.Password).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (p PostgresStorage) CreateUser(request models.LoginRequest) (models.User, error) {
	var user models.User
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id,login`
	err := p.Conn.QueryRow(context.Background(), query, request.Login, request.Password).Scan(&user.ID, &user.Login)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func NewPgStorage(ctx context.Context, dataBaseDSN string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(ctx, dataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool, %w", err)
	}
	return &PostgresStorage{Conn: pool}, nil
}

func (p PostgresStorage) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.Conn.Ping(ctx); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
