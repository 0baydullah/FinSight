package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/0baydullah/FinSight/internal/domain/auth"
	"github.com/0baydullah/FinSight/internal/requestctx"
	"github.com/0baydullah/FinSight/internal/response"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Authentication is not configured",
				"AUTH_CONFIGURATION_ERROR",
				nil,
			)
			return
		}

		header := r.Header.Get("Authorization")
		parts := strings.Fields(header)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(w)
			return
		}

		token, err := jwt.ParseWithClaims(
			parts[1],
			&auth.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				// Accept only the signing algorithm used by our login service.
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			},
		)

		if err != nil || !token.Valid {
			unauthorized(w)
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || claims.UserID == 0 {
			unauthorized(w)
			return
		}

		ctx := requestctx.WithUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func unauthorized(w http.ResponseWriter) {
	response.ErrorJSON(
		w,
		http.StatusUnauthorized,
		"Missing or invalid authentication token",
		"UNAUTHORIZED",
		nil,
	)
}
