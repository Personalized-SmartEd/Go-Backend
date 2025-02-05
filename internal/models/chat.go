package models

import "time"

type Chat struct {
	ChatID    string        `json:"chat_id" bson:"chat_id"`
	StudentID string        `json:"student_id" bson:"student_id"`
	Messages  []ChatMessage `json:"messages" bson:"messages"`
}

type ChatMessage struct {
	Content   string    `json:"content" bson:"content"`
	Sender    string    `json:"sender" bson:"sender"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
