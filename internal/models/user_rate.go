package models

import "time"

type UserRateCreate struct {
	UserID  int       `json:"user_id"`
	RateID  int       `json:"rate_id"`
	Lock    time.Time `json:"lock"`
	Deposit int       `json:"deposit"`
	Token   string    `json:"token"`
}

type UserRate struct {
	UserRateCreate
	ID          int       `json:"id"`
	Opened      time.Time `json:"opened"`
	EarnedPool  int       `json:"earned_pool"`
	OutcomePool int       `json:"outcome_pool"`
}

type UserRateClaim struct {
	Deposit int `json:"deposit"`
	Outcome int `json:"outcome"`
}
