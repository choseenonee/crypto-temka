package models

type WithdrawBase struct {
	Amount     float64     `json:"amount"`
	Token      string      `json:"token"`
	Status     string      `json:"status"`
	Properties interface{} `json:"properties"`
}

type WithdrawCreate struct {
	WithdrawBase
	UserID int `json:"user_id"`
}

type Withdraw struct {
	WithdrawCreate
	ID int `json:"id"`
}
