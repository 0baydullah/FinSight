package report

import (
	"context"
	"fmt"

	"github.com/0baydullah/FinSight/internal/shared"
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

func (r *reportRepository) GetDailyReport(
	ctx context.Context,
	filter ReportFilter,
) ([]ReportCategoryDto, []DailyCategoryData, error) {
	var categories []ReportCategoryDto

	// Get all active categories belonging to the user.
	err := r.db.WithContext(ctx).
		Table("categories").
		Select("id, name").
		Where("user_id = ?", filter.UserID).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Scan(&categories).Error
	if err != nil {
		return nil, nil, err
	}

	var dailyTotals []DailyCategoryData

	// Group expenses by date in the configured application timezone.
	dateExpression := fmt.Sprintf(
		"TO_CHAR(transaction_date AT TIME ZONE '%s', '%s')",
		shared.DefaultTimeZone,
		shared.PostgresDateFormat,
	)

	err = r.db.WithContext(ctx).
		Table("transactions").
		Select(
			dateExpression+" AS date, category_id, SUM(amount) AS amount",
		).
		Where("user_id = ?", filter.UserID).
		Where("type = ?", "expense").
		Where("transaction_date >= ?", filter.From).
		Where("transaction_date < ?", filter.To).
		Where("deleted_at IS NULL").
		Group(dateExpression + ", category_id").
		Order("date ASC").
		Order("category_id ASC").
		Scan(&dailyTotals).Error
	if err != nil {
		return nil, nil, err
	}

	return categories, dailyTotals, nil
}
