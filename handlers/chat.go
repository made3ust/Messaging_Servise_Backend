package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/models"
)

func CreateChat(c *gin.Context) {
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&chat)
	c.JSON(http.StatusCreated, chat)
}

func GetChats(c *gin.Context) {
	var chats []models.Chat
	config.DB.Find(&chats)
	c.JSON(http.StatusOK, chats)
}

func DeleteChat(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Chat{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
