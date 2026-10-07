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
