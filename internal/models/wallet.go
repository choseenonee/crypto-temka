package models

import "time"

type WalletUpdate struct {
	ID        int     `json:"id"`
	Token     string  `json:"token"`
	Deposit   float64 `json:"deposit"`
	IsOutcome bool    `json:"is_outcome,omitempty"`
	Outcome   float64 `json:"outcome"`
}

type Wallet struct {
	UserID int `json:"user_id"`
	WalletUpdate
}

type WalletInsertHistory struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"deposit"`
	WalletID  int       `json:"wallet_id"`
	TimeStamp time.Time `json:"timestamp"`
}
