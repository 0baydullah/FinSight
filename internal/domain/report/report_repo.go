package report

import (
	"context"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetSummary(
	ctx context.Context,
	filter ReportFilter,
) (*SummaryData, error) {
	var result SummaryData

	err := r.db.WithContext(ctx).
		Table("transactions").
		Select(`
            COALESCE(
                SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END),
                0
            ) AS total_income,
            COALESCE(
                SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END),
                0
            ) AS total_expenses
        `).
		Where("user_id = ?", filter.UserID).
		Where("transaction_date >= ?", filter.From).
		Where("transaction_date < ?", filter.To).
		Where("deleted_at IS NULL").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).
		Table("users").
		Select("balance").
		Where("id = ?", filter.UserID).
		Scan(&result.CurrentBalance).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
