package app

import (
	"net/http"

	"github.com/0baydullah/FinSight/internal/domain/user"
	"github.com/0baydullah/FinSight/internal/router"
	"gorm.io/gorm"
)

func SetupHandler(db *gorm.DB) http.Handler {
	// User dependencies
	userRepo := user.NewUserRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Register routes
	return router.Setup(userHandler)
}
