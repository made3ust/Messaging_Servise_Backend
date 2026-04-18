package main

import (
	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/handlers"
	"github.com/made3ust/Messaging_Servise_Backend/middleware"
)

func main() {
	config.ConnectDatabase()
	//config.DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUsers)
		protected.GET("/users/:id", handlers.GetUserByID)

		protected.POST("/chats", handlers.CreateChat)
		protected.GET("/chats", handlers.GetChats)

		protected.POST("/messages", handlers.SendMessage)
		protected.GET("/messages/chat/:chatId", handlers.GetMessages)
	}

	r.Run(":8080")
}

//migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/messenger?sslmode=disable" up
