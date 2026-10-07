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
