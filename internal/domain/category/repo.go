package category

import "context"

type Repository interface {
	Create(ctx context.Context, category *Category) error

	GetByID(ctx context.Context, id uint) (*Category, error)

	GetByUserID(ctx context.Context, userID uint) ([]Category, error)

	Update(ctx context.Context, category *Category) error

	Delete(ctx context.Context, category *Category) error
}
