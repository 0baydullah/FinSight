package transaction

import "context"

type Repository interface {
	Create(ctx context.Context, transaction *Transaction, balanceDelta int64) error

	GetByID(ctx context.Context, id uint) (*Transaction, error)

	List(ctx context.Context, filter TransactionFilter) (*TransactionPage, error)

	Update(ctx context.Context, transaction *Transaction, balanceDelta int64) error

	Delete(ctx context.Context, transaction *Transaction, balanceDelta int64) error
}