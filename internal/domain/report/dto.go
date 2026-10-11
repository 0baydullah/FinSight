package report

type SummaryResponseDto struct {
	From           string `json:"from"`
	To             string `json:"to"`
	TotalIncome    string `json:"total_income"`
	TotalExpenses  string `json:"total_expenses"`
	NetCashFlow    string `json:"net_cash_flow"`
	CurrentBalance string `json:"current_balance"`
}

type SummaryData struct {
	TotalIncome    int64
	TotalExpenses  int64
	CurrentBalance int64
}

type DailyReportResponseDto struct {
	From           string                `json:"from"`
	To             string                `json:"to"`
	Categories     []ReportCategoryDto   `json:"categories"`
	Days           []DailyReportRowDto   `json:"days"`
	CategoryTotals []CategorySubtotalDto `json:"category_totals"`
	TotalExpenses  string                `json:"total_expenses"`
}

type ReportCategoryDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DailyReportRowDto struct {
	Date            string                   `json:"date"`
	CategoryAmounts []DailyCategoryAmountDto `json:"category_amounts"`
	Total           string                   `json:"total"`
}

type DailyCategoryAmountDto struct {
	CategoryID uint   `json:"category_id"`
	Amount     string `json:"amount"`
}

type CategorySubtotalDto struct {
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name"`
	Total      string `json:"total"`
}

// Internal repository result; amounts are stored in paisa.
type DailyCategoryData struct {
	Date       string
	CategoryID uint
	Amount     int64
}
