package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	var req CreateTransactionRequestDto

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

	transaction, err := h.service.Create(
		r.Context(),
		userID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrTransactionCategory):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid category",
				"INVALID_CATEGORY",
				nil,
			)

		case errors.Is(err, ErrTransactionTypeInvalid):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Transaction type must be income or expense",
				"INVALID_TRANSACTION_TYPE",
				nil,
			)

		case errors.Is(err, ErrTransactionAmountInvalid):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid transaction amount",
				"INVALID_AMOUNT",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to create transaction",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		"Transaction created successfully",
		transaction,
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
			err.Error(),
		)
		return
	}

	query := r.URL.Query()

	// Pagination
	page := 1
	if value := query.Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid page",
				"INVALID_PAGE",
				"page must be a positive integer",
			)
			return
		}
		page = parsed
	}

	limit := 10
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid limit",
				"INVALID_LIMIT",
				"limit must be a positive integer",
			)
			return
		}
		limit = parsed
	}

	if limit > 100 {
		limit = 100
	}

	filter := TransactionFilter{
		Type:   strings.ToLower(strings.TrimSpace(query.Get("type"))),
		Search: strings.TrimSpace(query.Get("search")),
		Page:   page,
		Limit:  limit,
	}

	// Multiple categories: ?category_ids=1,3,5
	if value := query.Get("category_ids"); value != "" {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)

			id, err := strconv.ParseUint(item, 10, 64)
			if err != nil || id == 0 {
				response.ErrorJSON(
					w,
					http.StatusBadRequest,
					"Invalid category IDs",
					"INVALID_CATEGORY_IDS",
					"category_ids must contain comma-separated positive integers",
				)
				return
			}

			filter.CategoryIDs = append(filter.CategoryIDs, uint(id))
		}
	}

	// Date filters use YYYY-MM-DD.
	if value := query.Get("from"); value != "" {
		from, err := time.Parse("2006-01-02", value)
		if err != nil {
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid from date",
				"INVALID_FROM_DATE",
				"from must use YYYY-MM-DD format",
			)
			return
		}

		filter.From = &from
	}

	if value := query.Get("to"); value != "" {
		to, err := time.Parse("2006-01-02", value)
		if err != nil {
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid to date",
				"INVALID_TO_DATE",
				"to must use YYYY-MM-DD format",
			)
			return
		}

		// Exclusive upper bound includes the entire selected day.
		to = to.AddDate(0, 0, 1)
		filter.To = &to
	}

	if filter.From != nil &&
		filter.To != nil &&
		!filter.From.Before(*filter.To) {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid date range",
			"INVALID_DATE_RANGE",
			"from must be on or before to",
		)
		return
	}

	transactions, pagination, err := h.service.GetAll(
		r.Context(),
		userID,
		filter,
	)
	if err != nil {
		if errors.Is(err, ErrTransactionTypeInvalid) ||
			errors.Is(err, ErrTransactionCategory) {
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid filter",
				"INVALID_FILTER",
				err.Error(),
			)
			return
		}

		response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve transactions",
			"TRANSACTION_LIST_FAILED",
			err.Error(),
		)
		return
	}

	response.JSONWithPagination(
		w,
		http.StatusOK,
		"Transactions retrieved successfully",
		transactions,
		pagination,
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

	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid transaction ID",
			"INVALID_TRANSACTION_ID",
			nil,
		)
		return
	}

	transaction, err := h.service.GetByID(
		r.Context(),
		userID,
		uint(id),
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrTransactionNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Transaction not found",
				"TRANSACTION_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedTransaction):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this transaction",
				"TRANSACTION_FORBIDDEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to retrieve transaction",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Transaction retrieved successfully",
		transaction,
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

	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid transaction ID",
			"INVALID_TRANSACTION_ID",
			nil,
		)
		return
	}

	var req UpdateTransactionRequestDto

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

	transaction, err := h.service.Update(
		r.Context(),
		userID,
		uint(id),
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrTransactionNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Transaction not found",
				"TRANSACTION_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedTransaction):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this transaction",
				"TRANSACTION_FORBIDDEN",
				nil,
			)

		case errors.Is(err, ErrTransactionCategory):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid category",
				"INVALID_CATEGORY",
				nil,
			)

		case errors.Is(err, ErrTransactionTypeInvalid):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Transaction type must be income or expense",
				"INVALID_TRANSACTION_TYPE",
				nil,
			)

		case errors.Is(err, ErrTransactionAmountInvalid):
			response.ErrorJSON(
				w,
				http.StatusBadRequest,
				"Invalid transaction amount",
				"INVALID_AMOUNT",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to update transaction",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Transaction updated successfully",
		transaction,
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

	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"Invalid transaction ID",
			"INVALID_TRANSACTION_ID",
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
		case errors.Is(err, ErrTransactionNotFound):
			response.ErrorJSON(
				w,
				http.StatusNotFound,
				"Transaction not found",
				"TRANSACTION_NOT_FOUND",
				nil,
			)

		case errors.Is(err, ErrUnauthorizedTransaction):
			response.ErrorJSON(
				w,
				http.StatusForbidden,
				"You do not have access to this transaction",
				"TRANSACTION_FORBIDDEN",
				nil,
			)

		default:
			response.ErrorJSON(
				w,
				http.StatusInternalServerError,
				"Failed to delete transaction",
				"INTERNAL_ERROR",
				nil,
			)
		}

		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Transaction deleted successfully",
		nil,
	)
}
