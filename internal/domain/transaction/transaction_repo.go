package transaction

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Create(
	ctx context.Context,
	transaction *Transaction,
	balanceDelta int64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return tx.Table("users").
			Where("id = ?", transaction.UserID).
			UpdateColumn(
				"balance",
				gorm.Expr("balance + ?", balanceDelta),
			).
			Error
	})
}

func (r *TransactionRepository) GetByID(
	ctx context.Context,
	id uint,
) (*Transaction, error) {
	var transaction Transaction

	err := r.db.WithContext(ctx).
		First(&transaction, id).
		Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetByUserID(
	ctx context.Context,
	userID uint,
) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("transaction_date DESC").
		Find(&transactions).
		Error

	return transactions, err
}

func (r *TransactionRepository) Update(
	ctx context.Context,
	transaction *Transaction,
	balanceDelta int64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(transaction).Error; err != nil {
			return err
		}

		return tx.Table("users").
			Where("id = ?", transaction.UserID).
			UpdateColumn(
				"balance",
				gorm.Expr("balance + ?", balanceDelta),
			).
			Error
	})
}

func (r *TransactionRepository) Delete(
	ctx context.Context,
	transaction *Transaction,
	balanceDelta int64,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(transaction).Error; err != nil {
			return err
		}

		return tx.Table("users").
			Where("id = ?", transaction.UserID).
			UpdateColumn(
				"balance",
				gorm.Expr("balance + ?", balanceDelta),
			).
			Error
	})
}

func (r *TransactionRepository) List(
	ctx context.Context,
	filter TransactionFilter,
) (*TransactionPage, error) {
	query := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("user_id = ?", filter.UserID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if len(filter.CategoryIDs) > 0 {
		query = query.Where("category_id IN ?", filter.CategoryIDs)
	}

	if filter.From != nil {
		query = query.Where("transaction_date >= ?", *filter.From)
	}

	if filter.To != nil {
		query = query.Where("transaction_date < ?", *filter.To)
	}

	if filter.Search != "" {
		query = query.Where(
			"description ILIKE ?",
			"%"+filter.Search+"%",
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (filter.Page - 1) * filter.Limit

	var items []Transaction
	if err := query.
		Order("transaction_date DESC").
		Order("id DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(
			(total + int64(filter.Limit) - 1) /
				int64(filter.Limit),
		)
	}

	return &TransactionPage{
		Items:      items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}
