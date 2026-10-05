package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0baydullah/FinSight/internal/requestctx"
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

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := requestctx.GetUserID(r.Context())
	if err != nil {
		response.ErrorJSON(
			w,
			http.StatusUnauthorized,
			"Authentication required",
			"UNAUTHORIZED",
			nil,
		)
		return
	}

	profile, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"User not found",
				"USER_NOT_FOUND",
				nil,
			)
			return
		}

		response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"Internal server error",
			"INTERNAL_SERVER_ERROR",
			nil,
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Profile retrieved successfully",
		profile,
	)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, err := requestctx.GetUserID(r.Context())
	if err != nil {
		response.ErrorJSON(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
			"UNAUTHORIZED",
			nil,
		)
		return
	}

	var req UpdateUserRequestDto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
			"INVALID_REQUEST",
			nil,
		)
		return
	}

	user, err := h.service.UpdateMe(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNameRequired):
			response.ErrorJSON(w, http.StatusBadRequest,
				"Name is required", "NAME_REQUIRED", nil)

		case errors.Is(err, ErrUsernameRequired):
			response.ErrorJSON(w, http.StatusBadRequest,
				"Username is required", "USERNAME_REQUIRED", nil)

		case errors.Is(err, ErrEmailRequired):
			response.ErrorJSON(w, http.StatusBadRequest,
				"Email is required", "EMAIL_REQUIRED", nil)

		case errors.Is(err, ErrInvalidEmail):
			response.ErrorJSON(w, http.StatusBadRequest,
				"Invalid email address", "INVALID_EMAIL", nil)

		case errors.Is(err, ErrUsernameTaken):
			response.ErrorJSON(w, http.StatusConflict,
				"Username is already taken", "USERNAME_TAKEN", nil)

		case errors.Is(err, ErrEmailTaken):
			response.ErrorJSON(w, http.StatusConflict,
				"Email is already taken", "EMAIL_TAKEN", nil)

		case errors.Is(err, ErrPhoneTaken):
			response.ErrorJSON(w, http.StatusConflict,
				"Phone number is already taken", "PHONE_TAKEN", nil)

		case errors.Is(err, ErrUserNotFound):
			response.ErrorJSON(w, http.StatusNotFound,
				"User not found", "USER_NOT_FOUND", nil)

		default:
			response.ErrorJSON(w, http.StatusInternalServerError,
				"Failed to update user", "INTERNAL_ERROR", nil)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"User updated successfully",
		user,
	)
}
