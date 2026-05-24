package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/models"
	"github.com/made3ust/Messaging_Servise_Backend/utils"
)

func SendMessage(c *gin.Context) {
	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")

	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	msg.UserID = user.ID

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error: " + err.Error()})
		return
	}

	go utils.CallNotificationService(fmt.Sprintf("%v", username), msg.Content)

	c.JSON(http.StatusCreated, msg)
}

func GetMessages(c *gin.Context) {
	chatID := c.Param("chatId")
	var msgs []models.Message
	config.DB.Where("chat_id = ?", chatID).Find(&msgs)
	c.JSON(http.StatusOK, msgs)
}

func GetMessageByID(c *gin.Context) {
	id := c.Param("id")
	var msg models.Message
	if err := config.DB.First(&msg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func UpdateMessage(c *gin.Context) {
	id := c.Param("id")
	var msg models.Message
	if err := config.DB.First(&msg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg.Content = input.Content
	config.DB.Save(&msg)
	c.JSON(http.StatusOK, msg)
}

func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Message{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}

func SearchMessages(c *gin.Context) {
	query := c.Query("q")
	var msgs []models.Message
	config.DB.Where("content ILIKE ?", "%"+query+"%").Find(&msgs)
	c.JSON(http.StatusOK, msgs)
}
