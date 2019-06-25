package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterActions interface {
	Setup()
}

type MyRouter struct {
	*gin.Engine
}

func (router *MyRouter) Setup() {
	gin.SetMode(gin.DebugMode)
	router.Engine = gin.New()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authentication", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Authentication", "Authorization", "Content-Type"},
	}))

	public := router.Group("/")
	{
		public.POST("/Signup", Signup)
		public.POST("/Login", Login)
	}

	user := router.Group("/User")
	{
		user.POST("", CreateUser)
		user.GET("", GetUsers)
		user.GET("/:id_user", GetUser)
		user.PUT("/:id_user", UpdateUser)
		user.PUT("/:id_user/Enabled", ToggleUserEnabled)
	}

	//todo middleware
}
