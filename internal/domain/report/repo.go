package report

import "context"

type Repository interface {
	GetSummary(ctx context.Context, filter ReportFilter) (*SummaryData, error)
	GetDailyReport(ctx context.Context, filter ReportFilter) ([]ReportCategoryDto, []DailyCategoryData, error)
}
