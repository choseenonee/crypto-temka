package models

type WithdrawBase struct {
	Amount float64 `json:"amount"`
	//Token      string      `json:"token"`
	WalletID   int         `json:"wallet_id"`
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
