package models

type Wallet struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Token     string  `json:"token"`
	Deposit   float64 `json:"deposit"`
	IsOutcome bool    `json:"is_outcome,omitempty"`
}
