package models

type ReviewCreate struct {
	Title      string      `json:"tittle"`
	Text       string      `json:"text"`
	Properties interface{} `json:"properties"`
}

type Review struct {
	ReviewCreate
	ID int `json:"id"`
}

type MetricsSet struct {
	CurrentUsers  int `json:"current_users"`
	AlltimeIncome int `json:"alltime_income"`
	AlltimeOut    int `json:"alltime_out"`
}

type Metrics struct {
	MetricsSet
	IncomeSubOut int `json:"income_sub_out"`
}

type CaseCreate struct {
	Properties interface{} `json:"properties"`
}

type Case struct {
	CaseCreate
	ID int `json:"id"`
}

type Outcome struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
	Token  string  `json:"token"`
}
