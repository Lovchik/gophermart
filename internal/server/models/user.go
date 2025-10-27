package models

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenPair struct {
	Access  string
	Refresh string
}

type User struct {
	ID    int64  `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
}
