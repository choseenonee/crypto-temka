package models

type Wallet struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Token   string `json:"token"`
	Deposit int    `json:"deposit"`
}
