package models

import "time"

type Chat struct {
	StudentID     string    `bson:"student_id" validate:"required"`
	PreviousChats []Message `bson:"previous_chats"`
}

type Message struct {
	Prompt     string    `bson:"prompt"`
	Response   string    `bson:"response"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
