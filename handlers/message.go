package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/models"
)

func SendMessage(c *gin.Context) {
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&msg)
	c.JSON(http.StatusCreated, msg)
}

func GetMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	var msgs []models.Message
	config.DB.Where("chat_id = ?", chatID).Find(&msgs)
	c.JSON(http.StatusOK, msgs)
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Message{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
