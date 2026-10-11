package transaction

type CreateTransactionRequestDto struct {
	CategoryID      uint   `json:"category_id"`
	Type            string `json:"type"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
	TransactionDate string `json:"transaction_date"`
}

type UpdateTransactionRequestDto struct {
	CategoryID      *uint   `json:"category_id"`
	Type            *string `json:"type"`
	Amount          *string `json:"amount"`
	Description     *string `json:"description"`
	TransactionDate *string `json:"transaction_date"`
}

type TransactionResponseDto struct {
	ID              uint   `json:"id"`
	CategoryID      uint   `json:"category_id"`
	Type            string `json:"type"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
	TransactionDate string `json:"transaction_date"`
	CreatedAt       string `json:"created_at"`
}

type TransactionListResponseDto struct {
    Items      []TransactionResponseDto `json:"items"`
    Page       int                      `json:"page"`
    Limit      int                      `json:"limit"`
    Total      int64                    `json:"total"`
    TotalPages int                      `json:"total_pages"`
}

type TransactionPaginationDto struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}