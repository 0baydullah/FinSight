package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0baydullah/FinSight/internal/response"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequestDto

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

	user, err := h.service.CreateUser(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrNameRequired):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "NAME_REQUIRED", nil)

		case errors.Is(err, ErrUsernameRequired):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "USERNAME_REQUIRED", nil)

		case errors.Is(err, ErrEmailRequired):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "EMAIL_REQUIRED", nil)

		case errors.Is(err, ErrPasswordRequired):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "PASSWORD_REQUIRED", nil)

		case errors.Is(err, ErrInvalidEmail):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "INVALID_EMAIL", nil)

		case errors.Is(err, ErrWeakPassword):
			response.ErrorJSON(w, http.StatusBadRequest, err.Error(), "WEAK_PASSWORD", nil)

		case errors.Is(err, ErrUsernameTaken):
			response.ErrorJSON(w, http.StatusConflict, err.Error(), "USERNAME_TAKEN", nil)

		case errors.Is(err, ErrEmailTaken):
			response.ErrorJSON(w, http.StatusConflict, err.Error(), "EMAIL_TAKEN", nil)

		case errors.Is(err, ErrPhoneTaken):
			response.ErrorJSON(w, http.StatusConflict, err.Error(), "PHONE_TAKEN", nil)

		case errors.Is(err, gorm.ErrRecordNotFound):
			response.ErrorJSON(w, http.StatusNotFound, "User not found", "USER_NOT_FOUND", nil)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Internal server error",
				"INTERNAL_SERVER_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		"User created successfully",
		user,
	)
}
