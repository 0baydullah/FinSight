package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0baydullah/FinSight/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
			"INVALID_REQUEST_BODY",
			nil,
		)
		return
	}

	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrIdentifierRequired):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "IDENTIFIER_REQUIRED", nil,
			)

		case errors.Is(err, ErrPasswordRequired):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "PASSWORD_REQUIRED", nil,
			)

		case errors.Is(err, ErrInvalidCredentials):
			response.ErrorJSON(
				w, http.StatusUnauthorized,
				"Invalid username/email or password",
				"INVALID_CREDENTIALS", nil,
			)

		case errors.Is(err, ErrJWTConfiguration):
			response.ErrorJSON(
				w, http.StatusInternalServerError,
				"Authentication is not configured",
				"AUTH_CONFIGURATION_ERROR", nil,
			)

		default:
			response.ErrorJSON(
				w, http.StatusInternalServerError,
				"Internal server error",
				"INTERNAL_SERVER_ERROR", nil,
			)
		}
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Login successful",
		result,
	)
}
