package app

import (
	"log"
	"net/http"
	"os"

	"github.com/0baydullah/FinSight/internal/database"
)

func Run() error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")

	handler := SetupHandler(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("FinSight server running on port %s", port)

	return server.ListenAndServe()
}
