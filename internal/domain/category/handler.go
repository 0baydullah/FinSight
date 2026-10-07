package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/0baydullah/FinSight/internal/requestctx"
	"github.com/0baydullah/FinSight/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req CreateCategoryRequestDto

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

	category, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNameRequired):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Category name is required",
				"CATEGORY_NAME_REQUIRED",
				nil,
			)

		case errors.Is(err, ErrCategoryNameTaken):
			response.ErrorJSON(
				w,
				http.StatusConflict,
				"Category already exists",
				"CATEGORY_NAME_TAKEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to create category",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		"Category created successfully",
		category,
	)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	categories, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve categories",
			"INTERNAL_ERROR",
			nil,
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Categories retrieved successfully",
		categories,
	)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid category ID",
			"INVALID_CATEGORY_ID",
			nil,
		)
		return
	}

	category, err := h.service.GetByID(
		r.Context(),
		userID,
		uint(id),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Category not found",
				"CATEGORY_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedCategory):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this category",
				"CATEGORY_FORBIDDEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to retrieve category",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Category retrieved successfully",
		category,
	)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid category ID",
			"INVALID_CATEGORY_ID",
			nil,
		)
		return
	}

	var req UpdateCategoryRequestDto

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

	category, err := h.service.Update(
		r.Context(),
		userID,
		uint(id),
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Category not found",
				"CATEGORY_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedCategory):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this category",
				"CATEGORY_FORBIDDEN",
				nil,
			)

		case errors.Is(err, ErrCategoryNameRequired):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Category name is required",
				"CATEGORY_NAME_REQUIRED",
				nil,
			)

		case errors.Is(err, ErrCategoryNameTaken):
			response.ErrorJSON(
				w,
				http.StatusConflict,
				"Category already exists",
				"CATEGORY_NAME_TAKEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to update category",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Category updated successfully",
		category,
	)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid category ID",
			"INVALID_CATEGORY_ID",
			nil,
		)
		return
	}

	err = h.service.Delete(
		r.Context(),
		userID,
		uint(id),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Category not found",
				"CATEGORY_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedCategory):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this category",
				"CATEGORY_FORBIDDEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to delete category",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Category deleted successfully",
		nil,
	)
}
