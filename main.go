// main.go
package main

import (
	"github.com/Andrewalifb/alpha-pos-system-email-service/service"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	emailService := service.NewEmailService()
	emailService.StartConsuming()
}
