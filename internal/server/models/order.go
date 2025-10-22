package models

import "time"

type Order struct {
	ID         int64     `json:"-" gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	UserID     string    `json:"-"`
	Status     string    `json:"status"`
	Accurual   *string   `json:"accurual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	Number     string    `json:"number"`
}
