package router

import (
	"github.com/rs/cors"
	"gitlab.com/horo-go/horo"
	"gitlab.com/horo-go/pathvariable"
	"gitlab.com/horo-go/recovery"
	"gitlab.com/lattetalk/lattetalk/auth"
	"gitlab.com/lattetalk/lattetalk/config"
	"gitlab.com/lattetalk/lattetalk/controller"
)

func RegisterRoutes(router horo.Router) {
	rs := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
	})
	router.Use(rs.Handler)

	router.POST("/authorize", controller.Authorize)

	router.Use(auth.Auth0.Handler, auth.Auth0.Authenticated)

	rc := recovery.NewRecovery()
	router.Use(rc.Handler)

	router.POST("/logout", controller.LogOut)

	router.GET("/profile", controller.GetProfile)
	router.GET("/conversations", controller.GetConversations)
	router.GET("/friends", controller.GetFriends)
	router.GET("/friends/search", controller.SearchNewFriendsWithName)
	router.
		With(pathvariable.CheckInt64Min("id", 1)).
		POST("/friends/add/:id", controller.AddFriend)
	router.
		With(pathvariable.CheckInt64Min("id", 1)).
		GET("/conversations/:id/messages", controller.GetConversationMessage)
}
