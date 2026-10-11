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
			"Unauthorized",
			"UNAUTHORIZED",
			err.Error(),
		)
		return
	}

	summary, err := h.service.GetSummary(
		r.Context(),
		userID,
		r.URL.Query().Get("from"),
		r.URL.Query().Get("to"),
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidDate):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid date",
				"INVALID_DATE",
				err.Error(),
			)
		case errors.Is(err, ErrInvalidRange):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid date range",
				"INVALID_DATE_RANGE",
				err.Error(),
			)
		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to retrieve financial summary",
				"SUMMARY_FAILED",
				err.Error(),
			)
		}
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Financial summary retrieved successfully",
		summary,
	)
}
