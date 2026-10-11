package shared

const (
	DefaultTimeZone = "Asia/Dhaka"

	// Change this one value to change the app's display format.
	DateTimeFormat = "02 Jan 2006, 03:04 PM"

	// Internal/API date format, used for parsing date filters.
	DateFormat = "2006-01-02"

	// Database-specific format for PostgreSQL TO_CHAR.
	PostgresDateFormat = "YYYY-MM-DD"
)
