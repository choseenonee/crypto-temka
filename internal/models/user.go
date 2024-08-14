package models

type UserBase struct {
	Email       string      `json:"email"`
	PhoneNumber string      `json:"phone_number"`
	ReferID     *int        `json:"refer_id"`
	Properties  interface{} `json:"properties"`
}

type UserCreate struct {
	UserBase
	Password string `json:"password"`
}

type User struct {
	UserBase
	ID     int    `json:"id"`
	Status string `json:"status"`
}
