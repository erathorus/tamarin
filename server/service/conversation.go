package service

import (
	"database/sql"

	"gitlab.com/lattetalk/lattetalk/db"
	"gitlab.com/lattetalk/lattetalk/model"
)

// TODO: Do this inside a transaction when support delete
// Do not user to get user ids of other user when this query is
// anonymous
func GetConversationsByUserID(id int64) []model.Conversation {
	stmt, err := db.DB.Prepare(
		`SELECT id, created_at, updated_at
				FROM user_conversation JOIN conversation ON conversation_id=id
				WHERE user_id=$1`,
	)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}

	// TODO: change 20 to the max conversation query
	conversations := make([]model.Conversation, 0, 20)
	var conversation model.Conversation
	for rows.Next() {
		if err = rows.Scan(&conversation.ID, &conversation.CreatedAt, &conversation.UpdatedAt); err != nil {
			panic(err)
		}
		conversation.UserIDs = GetConversationUsersByID(conversation.ID)
		conversations = append(conversations, conversation)
	}
	return conversations
}

func GetConversationUsersByID(id int64) []int64 {
	stmt, err := db.DB.Prepare(
		`SELECT user_id
				FROM user_conversation
				WHERE conversation_id=$1`,
	)
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}

	userIDs := make([]int64, 0, 2)
	var userID int64
	for rows.Next() {
		if err = rows.Scan(&userID); err != nil {
			panic(err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

func ShouldGetConversationUsersByID(id int64) ([]int64, error) {
	stmt, err := db.DB.Prepare(
		`SELECT user_id
				FROM user_conversation
				WHERE conversation_id=$1`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, 0, 2)
	var userID int64
	for rows.Next() {
		if err = rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func GetConversationMessages(id int64, before int64) []model.Message {
	var rows *sql.Rows
	var err error
	if before == 0 {
		stmt, err := db.DB.Prepare(
			`SELECT id, user_id, created_at, content
					FROM message
					WHERE conversation_id=$1
					ORDER BY id DESC
					LIMIT 30`,
		)
		if err != nil {
			panic(err)
		}
		rows, err = stmt.Query(id)
	} else {
		stmt, err := db.DB.Prepare(
			`SELECT id, user_id, created_at, content
					FROM message
					WHERE conversation_id=$1 AND id < $2
					ORDER BY id DESC
					LIMIT 30`,
		)
		if err != nil {
			panic(err)
		}
		rows, err = stmt.Query(id, before)
	}

	if err != nil {
		panic(err)
	}

	messages := make([]model.Message, 0)
	var message model.Message
	message.ConversationID = id
	for rows.Next() {
		if err = rows.Scan(&message.ID, &message.UserID, &message.CreatedAt, &message.Content); err != nil {
			panic(err)
		}
		messages = append(messages, message)
	}

	return messages
}
