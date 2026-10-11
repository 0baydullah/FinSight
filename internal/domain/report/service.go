package report

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0baydullah/FinSight/internal/shared"
)

var (
	ErrInvalidDate  = errors.New("invalid date format; expected YYYY-MM-DD")
	ErrInvalidRange = errors.New("from date must be on or before to date")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetSummary(
	ctx context.Context,
	userID uint,
	fromDate string,
	toDate string,
) (*SummaryResponseDto, error) {
	location, err := time.LoadLocation(shared.DefaultTimeZone)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(location)

	// Default date range: first day of this month through today.
	defaultFrom := time.Date(
		today.Year(),
		today.Month(),
		1,
		0, 0, 0, 0,
		location,
	)

	defaultTo := time.Date(
		today.Year(),
		today.Month(),
		today.Day(),
		0, 0, 0, 0,
		location,
	)

	from := defaultFrom
	to := defaultTo

	if fromDate != "" {
		from, err = shared.ParseDate(fromDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	if toDate != "" {
		to, err = shared.ParseDate(toDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	if from.After(to) {
		return nil, ErrInvalidRange
	}

	// Use an exclusive upper bound to include the entire selected end date.
	toExclusive := to.AddDate(0, 0, 1)

	data, err := s.repo.GetSummary(ctx, ReportFilter{
		UserID: userID,
		From:   from,
		To:     toExclusive,
	})
	if err != nil {
		return nil, err
	}

	netCashFlow := data.TotalIncome - data.TotalExpenses

	return &SummaryResponseDto{
		From:           shared.FormatDate(from),
		To:             shared.FormatDate(to),
		TotalIncome:    formatMoney(data.TotalIncome),
		TotalExpenses:  formatMoney(data.TotalExpenses),
		NetCashFlow:    formatMoney(netCashFlow),
		CurrentBalance: formatMoney(data.CurrentBalance),
	}, nil
}

func formatMoney(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return fmt.Sprintf("%s%d.%02d", sign, amount/100, amount%100)
}

// FormatDateTime formats a timestamp using the shared short date-time layout.
func FormatDateTime(value time.Time) string {
	location, err := time.LoadLocation(shared.DefaultTimeZone)
	if err != nil {
		return value.Format(shared.DateTimeFormat)
	}

	return value.In(location).Format(shared.DateTimeFormat)
}

func (s *Service) GetDailyReport(
	ctx context.Context,
	userID uint,
	fromDate string,
	toDate string,
) (*DailyReportResponseDto, error) {
	location, err := time.LoadLocation(shared.DefaultTimeZone)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(location)

	// Default range: first day of this month through today.
	from := time.Date(
		today.Year(), today.Month(), 1,
		0, 0, 0, 0, location,
	)
	to := time.Date(
		today.Year(), today.Month(), today.Day(),
		0, 0, 0, 0, location,
	)

	if fromDate != "" {
		from, err = shared.ParseDate(fromDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	if toDate != "" {
		to, err = shared.ParseDate(toDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	if from.After(to) {
		return nil, ErrInvalidRange
	}

	// The repository uses an exclusive upper boundary.
	toExclusive := to.AddDate(0, 0, 1)

	categories, dailyTotals, err := s.repo.GetDailyReport(
		ctx,
		ReportFilter{
			UserID: userID,
			From:   from,
			To:     toExclusive,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create a lookup for daily category amounts.
	amountsByDate := make(map[string]map[uint]int64)

	for _, item := range dailyTotals {
		if amountsByDate[item.Date] == nil {
			amountsByDate[item.Date] = make(map[uint]int64)
		}

		amountsByDate[item.Date][item.CategoryID] = item.Amount
	}

	// Prepare category subtotals.
	categoryTotals := make(map[uint]int64)
	days := make([]DailyReportRowDto, 0)
	var grandTotal int64

	// Generate every date, including dates with no transactions.
	for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
		dateString := day.Format(shared.DateFormat)
		row := DailyReportRowDto{
			Date:            dateString,
			CategoryAmounts: make([]DailyCategoryAmountDto, 0, len(categories)),
		}

		var dailyTotal int64

		for _, category := range categories {
			amount := amountsByDate[dateString][category.ID]

			row.CategoryAmounts = append(
				row.CategoryAmounts,
				DailyCategoryAmountDto{
					CategoryID: category.ID,
					Amount:     formatMoney(amount),
				},
			)

			dailyTotal += amount
			categoryTotals[category.ID] += amount
		}

		row.Total = formatMoney(dailyTotal)
		grandTotal += dailyTotal
		days = append(days, row)
	}

	// Build category subtotal rows in the same order as the columns.
	subtotals := make([]CategorySubtotalDto, 0, len(categories))

	for _, category := range categories {
		subtotals = append(subtotals, CategorySubtotalDto{
			CategoryID: category.ID,
			Name:       category.Name,
			Total:      formatMoney(categoryTotals[category.ID]),
		})
	}

	return &DailyReportResponseDto{
		From:           from.Format(shared.DateFormat),
		To:             to.Format(shared.DateFormat),
		Categories:     categories,
		Days:           days,
		CategoryTotals: subtotals,
		TotalExpenses:  formatMoney(grandTotal),
	}, nil
}
