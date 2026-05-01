package utils

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallNotificationService(username string, text string) {
	client := resty.New()

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Printf("Resty: Sending data to notifdication service...\n")
		return nil
	})

	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"user":    username,
			"message": text,
		}).
		Post("http://localhost:8081/send-notification")

	if err != nil {
		fmt.Println("Resty Error:", err)
	}
}
