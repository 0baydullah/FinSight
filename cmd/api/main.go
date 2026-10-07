package main

import (
	"log"

	"github.com/0baydullah/FinSight/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
