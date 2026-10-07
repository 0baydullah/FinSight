package transaction

import "context"

type Repository interface {
	Create(ctx context.Context, transaction *Transaction, balanceDelta int64) error

	GetByID(ctx context.Context, id uint) (*Transaction, error)

	GetByUserID(ctx context.Context, userID uint) ([]Transaction, error)

	Update(ctx context.Context, transaction *Transaction, balanceDelta int64) error

	Delete(ctx context.Context, transaction *Transaction, balanceDelta int64) error
}
