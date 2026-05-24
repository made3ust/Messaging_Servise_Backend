package main

import (
	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/handlers"
	"github.com/made3ust/Messaging_Servise_Backend/middleware"
	"github.com/made3ust/Messaging_Servise_Backend/models"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{

		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/users", handlers.GetUsers)
		protected.GET("/users/:id", handlers.GetUserByID)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)
		protected.GET("/search/users", handlers.SearchUsers)

		protected.POST("/chats", handlers.CreateChat)
		protected.GET("/chats", handlers.GetChats)
		protected.GET("/chats/:id", handlers.GetChatByID)
		protected.PUT("/chats/:id", handlers.UpdateChat)
		protected.DELETE("/chats/:id", handlers.DeleteChat)
		protected.GET("/search/chats", handlers.SearchChats)

		protected.POST("/messages", handlers.SendMessage)
		protected.GET("/messages/chat/:chatId", handlers.GetMessages)
		protected.GET("/messages/:id", handlers.GetMessageByID)
		protected.PUT("/messages/:id", handlers.UpdateMessage)
		protected.DELETE("/messages/:id", handlers.DeleteMessage)
		protected.GET("/search/messages", handlers.SearchMessages)
	}

	r.Run(":8080")
}
