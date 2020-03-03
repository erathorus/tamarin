package controller

import (
	"net/http"

	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/service"
	"gitlab.com/lattetalk/lattetalk/ws"
	"gitlab.com/lattetalk/lattetalk/ws/data"
)

func GetFriends(c *horo.Context) {
	friends := service.GetUserFriendsByID(c.SubjectInt64())
	c.JSON(http.StatusOK, friends)
}

func SearchNewFriendsWithName(c *horo.Context) {
	name := c.Request.URL.Query().Get("name")
	if len(name) < 3 {
		c.BadRequest()
		return
	}
	friends := service.GetUserNewFriendsWithSimilarName(c.SubjectInt64(), name)
	c.JSON(http.StatusOK, friends)
}

func AddFriend(c *horo.Context) {
	friendID := c.ParamInt64("id")
	conversation, err := service.ShouldAddNewFriend(c.SubjectInt64(), friendID)
	if err != nil {
		http.Error(c.Writer, "internal error", http.StatusInternalServerError)
		return
	}
	client := ws.Hub.GetClient(friendID)
	client.Broadcast(data.ResponseNewFriend, service.GetUserByID(c.SubjectInt64()))
	client.Broadcast(data.ResponseNewConversation, conversation)
	c.JSON(http.StatusOK, conversation)
}
