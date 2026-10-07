package router

import (
	"encoding/json"
	"net/http"

	authDomain "github.com/0baydullah/FinSight/internal/domain/auth"
	categoryDomain "github.com/0baydullah/FinSight/internal/domain/category"
	userDomain "github.com/0baydullah/FinSight/internal/domain/user"
	"github.com/0baydullah/FinSight/internal/middleware"
)

func Setup(
	userHandler *userDomain.Handler,
	authHandler *authDomain.Handler,
	categoryHandler *categoryDomain.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// auth and user routes
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.Handle("GET /users/me", middleware.Auth(http.HandlerFunc(userHandler.GetMe)))
	mux.Handle("PATCH /users/me", middleware.Auth(http.HandlerFunc(userHandler.UpdateMe)))

	// category routes
	mux.Handle(
		"POST /categories",
		middleware.Auth(http.HandlerFunc(categoryHandler.Create)),
	)
	mux.Handle(
		"GET /categories",
		middleware.Auth(http.HandlerFunc(categoryHandler.GetAll)),
	)
	mux.Handle(
		"GET /categories/{id}",
		middleware.Auth(http.HandlerFunc(categoryHandler.GetByID)),
	)
	mux.Handle(
		"PATCH /categories/{id}",
		middleware.Auth(http.HandlerFunc(categoryHandler.Update)),
	)
	mux.Handle(
		"DELETE /categories/{id}",
		middleware.Auth(http.HandlerFunc(categoryHandler.Delete)),
	)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
