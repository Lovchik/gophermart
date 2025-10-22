package models

import "time"

type WithdrawalOrders struct {
	ID          int64     `json:"-" gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	UserID      string    `json:"-"`
	Status      string    `json:"-"`
	Accurual    *string   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
	Order       string    `json:"order"`
}

type CreateWithdrawalOrder struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required"`
}
