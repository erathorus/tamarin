package model

import (
	"time"
)

type Message struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"userId"`
	ConversationID int64     `json:"conversationId"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}
