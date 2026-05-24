package utils

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

func CallNotificationService(username string, text string) {
	notifURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notifURL == "" {
		notifURL = "http://localhost:8081"
	}

	client := resty.New()

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Printf("Resty: Sending data to notification service at %s...\n", notifURL)
		return nil
	})

	endpoint := fmt.Sprintf("%s/send-notification", notifURL)

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"user":    username,
			"message": text,
		}).
		Post(endpoint)

	if err != nil {
		fmt.Println("Resty Error:", err)
	}
}
