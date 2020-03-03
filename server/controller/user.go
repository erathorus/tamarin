package controller

import (
	"net/http"

	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/service"
)

func GetUserByID(c *horo.Context) {
	user := service.GetUserByOtherUserID(c.ParamInt64("id"), c.SubjectInt64())
	c.JSON(http.StatusOK, user)
}
