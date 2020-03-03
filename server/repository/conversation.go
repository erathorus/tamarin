package repository

import (
	"github.com/hashicorp/golang-lru"
	"gitlab.com/lattetalk/lattetalk/model"
	"gitlab.com/lattetalk/lattetalk/service"
)

const (
	conversationsCacheSize = 1 << 20
)

type conversations struct {
	conversations  *lru.Cache
	nConversations int
}

func (c *conversations) GetConversation(id int64) (*model.Conversation, error) {
	conversation, ok := c.conversations.Get(id)
	if !ok {
		userIDs, err := service.ShouldGetConversationUsersByID(id)
		if err != nil {
			return nil, err
		}
		conversation = &model.Conversation{
			ID:      id,
			UserIDs: userIDs,
		}
		c.conversations.Add(id, conversation)
	}
	return conversation.(*model.Conversation), nil
}

var Conversations *conversations

func init() {
	Conversations = new(conversations)

	var err error
	if Conversations.conversations, err = lru.New(conversationsCacheSize); err != nil {
		panic(err)
	}
}
