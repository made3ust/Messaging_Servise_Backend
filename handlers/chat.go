package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/models"
)

func CreateChat(c *gin.Context) {
	var input struct {
		Name    string `json:"name"`
		IsGroup bool   `json:"is_group"`
		UserIDs []uint `json:"user_ids"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")
	var creator models.User
	if err := config.DB.Where("username = ?", username).First(&creator).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creator not found"})
		return
	}

	var invitedUsers []models.User
	if len(input.UserIDs) > 0 {
		config.DB.Where("id IN ?", input.UserIDs).Find(&invitedUsers)
	}

	invitedUsers = append(invitedUsers, creator)

	chat := models.Chat{
		Name:    input.Name,
		IsGroup: input.IsGroup,
		Users:   invitedUsers,
	}

	if err := config.DB.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

func GetChats(c *gin.Context) {
	username, _ := c.Get("username")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var chats []models.Chat
	err := config.DB.Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
		Where("user_chats.user_id = ?", user.ID).
		Preload("Users").
		Preload("Messages").
		Find(&chats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func GetChatByID(c *gin.Context) {
	id := c.Param("id")
	var chat models.Chat
	if err := config.DB.Preload("Users").First(&chat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}
	c.JSON(http.StatusOK, chat)
}

func UpdateChat(c *gin.Context) {
	id := c.Param("id")
	var chat models.Chat
	if err := config.DB.First(&chat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&chat)
	c.JSON(http.StatusOK, chat)
}

func DeleteChat(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.Chat{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

func SearchChats(c *gin.Context) {
	query := c.Query("q")
	username, _ := c.Get("username")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var chats []models.Chat
	err := config.DB.Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
		Where("user_chats.user_id = ? AND chats.name ILIKE ?", user.ID, "%"+query+"%").
		Preload("Users").
		Preload("Messages").
		Find(&chats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}
