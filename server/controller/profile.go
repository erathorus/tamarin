package controller

import (
	"net/http"

	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/service"
)

func GetProfile(c *horo.Context) {
	user := service.GetUserByID(c.SubjectInt64())
	c.JSON(http.StatusOK, user)
}

