package service

import (
	"gitlab.com/lattetalk/lattetalk/db"
	"gitlab.com/lattetalk/lattetalk/model"
)

func ShouldCreateMessage(message *model.Message) (int64, error) {
	stmt, err := db.DB.Prepare(
		`INSERT INTO message (user_id, conversation_id, content, created_at)
				VALUES ($1, $2, $3, $4)
				RETURNING id`,
	)
	if err != nil {
		return -1, err
	}
	var id int64
	err = stmt.QueryRow(message.UserID, message.ConversationID, message.Content, message.CreatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
