package models

import (
	"time"
)

type UserRateCreate struct {
	UserID    int       `json:"user_id,omitempty"`
	RateID    int       `json:"rate_id"`
	Lock      time.Time `json:"lock"`
	Deposit   float64   `json:"deposit"`
	WalletID  int       `json:"wallet_id"`
	VoucherID *string   `json:"voucher_id"`
}

type UserRate struct {
	UserRateCreate
	ID          int       `json:"id"`
	Opened      time.Time `json:"opened"`
	EarnedPool  float64   `json:"earned_pool"`
	OutcomePool float64   `json:"outcome_pool"`
	Token       string    `json:"token"`
}

type UserRateAdmin struct {
	UserRate
	NextDayCharge *float64 `json:"next_day_charge"`
}
