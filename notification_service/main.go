package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/send-notification", func(c *gin.Context) {
		var req struct {
			User    string `json:"user"`
			Message string `json:"message"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf(" Sending notification to %s: %s\n", req.User, req.Message)

		c.JSON(200, gin.H{"status": "delivered"})
	})

	fmt.Println("Notification service started on :8081")
	r.Run(":8081")
}
