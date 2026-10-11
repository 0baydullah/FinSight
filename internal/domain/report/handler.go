package report

import (
	"errors"
	"net/http"

	"github.com/0baydullah/FinSight/internal/requestctx"
	"github.com/0baydullah/FinSight/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.GetSummary(
		r.Context(),
		userID,
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidDate):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "INVALID_DATE", nil,
			)
		case errors.Is(err, ErrInvalidRange):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "INVALID_DATE_RANGE", nil,
			)
		default:
			response.ErrorJSON(
				w, http.StatusInternalServerError,
				"Failed to generate summary", "INTERNAL_ERROR", nil,
			)
		}
		return
	}

	response.JSON(w, http.StatusOK, "Report summary retrieved successfully", result)
}

func (h *Handler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.GetDailyReport(
		r.Context(),
		userID,
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidDate):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "INVALID_DATE", nil,
			)
		case errors.Is(err, ErrInvalidRange):
			response.ErrorJSON(
				w, http.StatusBadRequest,
				err.Error(), "INVALID_DATE_RANGE", nil,
			)
		default:
			response.ErrorJSON(
				w, http.StatusInternalServerError,
				"Failed to generate daily report", "INTERNAL_ERROR", nil,
			)
		}
		return
	}

	response.JSON(w, http.StatusOK, "Daily report retrieved successfully", result)
}
