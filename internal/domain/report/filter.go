package report

import "time"

type ReportFilter struct {
	UserID uint
	From   time.Time
	To     time.Time
}
