package transaction

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
			nil,
		)
		return
	}

	transactions, err := h.service.GetAll(
		r.Context(),
		userID,
	)

	if err != nil {
		response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve transactions",
			"INTERNAL_ERROR",
			nil,
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Transactions retrieved successfully",
		transactions,
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
