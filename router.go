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
		public.POST("/SignUp", SignUp)
		public.POST("/LogIn", Login)
	}
}
