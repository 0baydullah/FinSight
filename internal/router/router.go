package router

import (
	"encoding/json"
	"net/http"

	authDomain "github.com/0baydullah/FinSight/internal/domain/auth"
	userDomain "github.com/0baydullah/FinSight/internal/domain/user"
)

func Setup(
	userHandler *userDomain.Handler,
	authHandler *authDomain.Handler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

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