package config

import (
	"fmt"
	"os" // Добавили для работы с переменными окружения

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=messenger port=5432 sslmode=disable", dbHost)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	fmt.Printf("Database connection established (Host: %s)\n", dbHost)
	DB = database
}
