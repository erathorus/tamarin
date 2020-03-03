package controller

import (
	"net/http"
	"strconv"

	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/repository"
	"gitlab.com/lattetalk/lattetalk/service"
)

func GetConversations(c *horo.Context) {
	conversations := service.GetConversationsByUserID(c.SubjectInt64())
	c.JSON(http.StatusOK, conversations)
}

func GetConversationMessage(c *horo.Context) {
	conversation, err := repository.Conversations.GetConversation(c.ParamInt64("id"))
	if err != nil {
		c.BadRequest()
		return
	}
	if !conversation.HasUser(c.SubjectInt64()) {
		c.BadRequest()
		return
	}
	var before int64
	beforeParam := c.Request.URL.Query().Get("before")
	if beforeParam == "" {
		before = 0
	} else {
		before, err = strconv.ParseInt(beforeParam, 10, 64)
		if err != nil {
			c.BadRequest()
			return
		}
	}
	c.JSON(http.StatusOK, service.GetConversationMessages(conversation.ID, before))
}
