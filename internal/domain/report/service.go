package report

import (
	"context"
	"errors"
	"fmt"
	"time"
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
	location, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		return nil, err
	}

	today := time.Now().In(location)
	defaultFrom := time.Date(
		today.Year(), today.Month(), 1, 0, 0, 0, 0, location,
	)
	defaultTo := today

	from := defaultFrom
	to := defaultTo

	if fromDate != "" {
		from, err = time.ParseInLocation("2006-01-02", fromDate, location)
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	if toDate != "" {
		to, err = time.ParseInLocation("2006-01-02", toDate, location)
		if err != nil {
			return nil, ErrInvalidDate
		}
	}

	// Compare calendar dates, not the time of day.
	if from.After(to) {
		return nil, ErrInvalidRange
	}

	// The repository uses an exclusive upper bound.
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
		From:           from.Format("2006-01-02"),
		To:             to.Format("2006-01-02"),
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
