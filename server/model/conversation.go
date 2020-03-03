package model

import (
	"time"
)

type Conversation struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserIDs   []int64   `json:"userIds"`
}

func (c *Conversation) HasUser(userID int64) bool {
	for _, id := range c.UserIDs {
		if userID == id {
			return true
		}
	}
	return false
}
