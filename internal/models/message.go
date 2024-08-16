package models

import "time"

type MessageCreate struct {
	UserID     int         `json:"user_id"`
	Properties interface{} `json:"properties"`
	Timestamp  time.Time   `json:"timestamp"`
}

type Message struct {
	MessageCreate
	ID     int  `json:"id"`
	IsRead bool `json:"is_read"`
}

//type UserMessage struct {
//	UserID   int       `json:"user_id"`
//	Messages []Message `json:"messages"`
//}
