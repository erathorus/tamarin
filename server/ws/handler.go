package ws

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"gitlab.com/lattetalk/lattetalk/model"
	"gitlab.com/lattetalk/lattetalk/repository"
	"gitlab.com/lattetalk/lattetalk/service"
	"gitlab.com/lattetalk/lattetalk/ws/data"
)

var (
	ErrConversationIsNotExist  = errors.New("conversations is not exist")
	ErrUserIsNotInConversation = errors.New("message sender is not in conversations")
)

func handleNewMessage(r *data.Request) error {
	message := new(model.Message)
	if err := json.Unmarshal(r.Data, message); err != nil {
		log.Print(err)
		return err
	}
	// Don't handle message with empty content
	if message.Content == "" {
		return nil
	}
	message.UserID = r.Subject()
	message.CreatedAt = time.Now()

	conversation, err := repository.Conversations.GetConversation(message.ConversationID)
	if err != nil {
		return ErrConversationIsNotExist
	}
	if !conversation.HasUser(message.UserID) {
		return ErrUserIsNotInConversation
	}

	id, err := service.ShouldCreateMessage(message)
	if err != nil {
		log.Print(err)
		return err
	}
	message.ID = id

	for _, userID := range conversation.UserIDs {
		client := Hub.GetClient(userID)
		if client == nil {
			continue
		}
		client.Broadcast(data.ResponseNewMessage, message)
	}

	return nil
}
