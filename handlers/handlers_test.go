package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/made3ust/Messaging_Servise_Backend/config"
	"github.com/made3ust/Messaging_Servise_Backend/models"
	"github.com/made3ust/Messaging_Servise_Backend/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})
	config.DB = db
}

func mockAuthMiddleware(username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("username", username)
		c.Next()
	}
}

func TestRegister_Success(t *testing.T) {
	setupTestDB()

	r := gin.New()
	r.POST("/register", Register)

	payload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["message"] != "User registered successfully" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}
}

func TestRegister_DuplicateUser(t *testing.T) {
	setupTestDB()

	hashedPassword, _ := utils.HashPassword("password123")
	existingUser := models.User{Username: "testuser", Email: "test@example.com", Password: hashedPassword}
	config.DB.Create(&existingUser)

	r := gin.New()
	r.POST("/register", Register)

	payload := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	setupTestDB()

	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{Username: "testuser", Email: "test@example.com", Password: hashedPassword}
	config.DB.Create(&user)

	r := gin.New()
	r.POST("/login", Login)

	payload := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["token"] == "" {
		t.Error("Expected token in response, got empty string")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	setupTestDB()

	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{Username: "testuser", Email: "test@example.com", Password: hashedPassword}
	config.DB.Create(&user)

	r := gin.New()
	r.POST("/login", Login)

	payload := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestCreateChat_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	r := gin.New()
	r.Use(mockAuthMiddleware("testuser"))
	r.POST("/chats", CreateChat)

	payload := map[string]interface{}{
		"name":     "General Chat",
		"is_group": true,
		"user_ids": []uint{},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/chats", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Chat
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Name != "General Chat" {
		t.Errorf("Expected chat name 'General Chat', got %s", response.Name)
	}
}

func TestGetChats_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	chat1 := models.Chat{Name: "Room 1", IsGroup: false, Users: []models.User{user}}
	chat2 := models.Chat{Name: "Room 2", IsGroup: true, Users: []models.User{user}}
	config.DB.Create(&chat1)
	config.DB.Create(&chat2)

	r := gin.New()
	r.Use(mockAuthMiddleware("testuser"))
	r.GET("/chats", GetChats)

	req, _ := http.NewRequest("GET", "/chats", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []models.Chat
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 2 {
		t.Errorf("Expected 2 chats, got %d", len(response))
	}
}

func TestSendMessage_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	chat := models.Chat{Name: "General", IsGroup: true, Users: []models.User{user}}
	config.DB.Create(&chat)

	r := gin.New()
	r.Use(mockAuthMiddleware("testuser"))
	r.POST("/messages", SendMessage)

	payload := map[string]interface{}{
		"content": "Hello World",
		"chat_id": chat.ID,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 211, got %d", w.Code)
	}

	var response models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Content != "Hello World" {
		t.Errorf("Expected message content 'Hello World', got %s", response.Content)
	}
}

func TestGetMessages_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	chat := models.Chat{Name: "General", IsGroup: true, Users: []models.User{user}}
	config.DB.Create(&chat)

	msg1 := models.Message{Content: "First", UserID: user.ID, ChatID: chat.ID}
	msg2 := models.Message{Content: "Second", UserID: user.ID, ChatID: chat.ID}
	config.DB.Create(&msg1)
	config.DB.Create(&msg2)

	r := gin.New()
	r.GET("/messages/chat/:chatId", GetMessages)

	req, _ := http.NewRequest("GET", "/messages/chat/"+strconv.Itoa(int(chat.ID)), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(response))
	}
}

func TestUpdateMessage_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	chat := models.Chat{Name: "General", IsGroup: true, Users: []models.User{user}}
	config.DB.Create(&chat)

	msg := models.Message{Content: "Original Content", UserID: user.ID, ChatID: chat.ID}
	config.DB.Create(&msg)

	r := gin.New()
	r.PUT("/messages/:id", UpdateMessage)

	payload := map[string]string{
		"content": "Updated Content",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/messages/"+strconv.Itoa(int(msg.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.Message
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Content != "Updated Content" {
		t.Errorf("Expected content 'Updated Content', got %s", response.Content)
	}
}

func TestDeleteMessage_Success(t *testing.T) {
	setupTestDB()

	user := models.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	config.DB.Create(&user)

	chat := models.Chat{Name: "General", IsGroup: true, Users: []models.User{user}}
	config.DB.Create(&chat)

	msg := models.Message{Content: "To Be Deleted", UserID: user.ID, ChatID: chat.ID}
	config.DB.Create(&msg)

	r := gin.New()
	r.DELETE("/messages/:id", DeleteMessage)

	req, _ := http.NewRequest("DELETE", "/messages/"+strconv.Itoa(int(msg.ID)), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var check models.Message
	err := config.DB.First(&check, msg.ID).Error
	if err == nil {
		t.Error("Expected message to be deleted from active query, but it was found")
	}
}
