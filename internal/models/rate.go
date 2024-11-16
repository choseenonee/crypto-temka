package models

type RateCreate struct {
	Title       string      `json:"title"`
	Profit      float64     `json:"profit"`
	MinLockDays int         `json:"min_lock_days"`
	Commission  int         `json:"commission"`
	Properties  interface{} `json:"properties"`
}

type Rate struct {
	RateCreate
	ID int `json:"id"`
}
