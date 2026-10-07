package app

import (
	"net/http"
	"os"
	"time"

	authDomain "github.com/0baydullah/FinSight/internal/domain/auth"
	categoryDomain "github.com/0baydullah/FinSight/internal/domain/category"
	transactionDomain "github.com/0baydullah/FinSight/internal/domain/transaction"
	userDomain "github.com/0baydullah/FinSight/internal/domain/user"
	"github.com/0baydullah/FinSight/internal/router"
	"gorm.io/gorm"
)

func SetupHandler(db *gorm.DB) http.Handler {
	// Repository
	userRepo := userDomain.NewUserRepository(db)

	// User service and handler
	userService := userDomain.NewService(userRepo)
	userHandler := userDomain.NewHandler(userService)

	// Authentication service and handler
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresIn := getJWTExpiration()

	authService := authDomain.NewService(
		userRepo,
		jwtSecret,
		jwtExpiresIn,
	)
	authHandler := authDomain.NewHandler(authService)

	// category dependencies
	categoryRepo := categoryDomain.NewCategoryRepository(db)
	categoryService := categoryDomain.NewService(categoryRepo)
	categoryHandler := categoryDomain.NewHandler(categoryService)

	// transaction dependencies
	transactionRepo := transactionDomain.NewTransactionRepository(db)

	transactionService := transactionDomain.NewService(
		transactionRepo,
		categoryRepo,
	)

	transactionHandler := transactionDomain.NewHandler(
		transactionService,
	)

	// Register routes
	return router.Setup(userHandler, authHandler, categoryHandler, transactionHandler)
}

func getJWTExpiration() time.Duration {
	value := os.Getenv("JWT_EXPIRES_IN")
	if value == "" {
		return 24 * time.Hour
	}

	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 24 * time.Hour
	}

	return duration
}
