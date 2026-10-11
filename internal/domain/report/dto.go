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
