package main

import (
	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/handlers"
	"github.com/made3ust/Messaging_Servise_Backend/models"
)

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", handlers.CreateUser)
		userRoutes.GET("/", handlers.GetUsers)
		userRoutes.GET("/:id", handlers.GetUserByID)
		userRoutes.PUT("/:id", handlers.UpdateUser)
	}

	chatRoutes := r.Group("/chats")
	{
		chatRoutes.POST("/", handlers.CreateChat)
		chatRoutes.GET("/", handlers.GetChats)
		chatRoutes.DELETE("/:id", handlers.DeleteChat)
	}

	msgRoutes := r.Group("/messages")
	{
		msgRoutes.POST("/", handlers.SendMessage)
		msgRoutes.GET("/chat/:chatId", handlers.GetMessages)
		msgRoutes.DELETE("/:id", handlers.DeleteMessage)
	}

	r.Run(":8080")
}
