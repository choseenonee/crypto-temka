package models

type UserBase struct {
	Email  string `json:"email"`
	Ref    int    `json:"id_ref"`
	Amount int    `json:"amount"`
}

type UserCreate struct {
	UserBase
	Password string `json:"password"`
}

type User struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Ref    int    `json:"id_ref"`
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPhoto struct {
	ID     int    `json:"userID"`
	Photo  []byte `json:"photo"`
	Status string `json:"status"`
}

type SetStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
