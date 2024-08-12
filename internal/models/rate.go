package models

import "time"

type RateCreate struct {
	Title       string      `json:"title"`
	Profit      int         `json:"profit"`
	MinLockDays int         `json:"min_lock_days"`
	Properties  interface{} `json:"properties"`
}

type Rate struct {
	RateCreate
	ID int `json:"id"`
}

type UserRateCreate struct {
	UserID  int       `json:"user_id"`
	RateID  int       `json:"rate_id"`
	Lock    time.Time `json:"lock"`
	Opened  time.Time `json:"opened"`
	Deposit int       `json:"deposit"`
	Token   string    `json:"token"`
}

type UserRate struct {
	UserRateCreate
	EarnedPool  int `json:"earned_pool"`
	OutcomePool int `json:"outcome_pool"`
}
