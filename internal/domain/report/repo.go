package report

import "context"

type Repository interface {
	GetSummary(ctx context.Context, filter ReportFilter) (*SummaryData, error)
}
