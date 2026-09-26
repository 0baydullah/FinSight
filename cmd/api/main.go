package main

import (
	"log"
	"net/http"

	"github.com/0baydullah/FinSight/internal/database"
	"github.com/0baydullah/FinSight/internal/domain/user"
	"github.com/0baydullah/FinSight/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found; using system environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL successfully!")

	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	handler := router.Setup(userHandler)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
