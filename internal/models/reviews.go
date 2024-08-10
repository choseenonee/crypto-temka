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
