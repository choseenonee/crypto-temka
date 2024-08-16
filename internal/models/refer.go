package models

type Refer struct {
	ID       int     `json:"id"`
	ParentID int     `json:"parent_id"`
	ChildID  int     `json:"child_id"`
	Amount   float64 `json:"amount"`
	Token    string  `json:"token"`
}
