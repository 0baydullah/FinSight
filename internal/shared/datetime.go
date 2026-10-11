package shared

import "time"

// FormatDateTime formats timestamps consistently across the application.
func FormatDateTime(value time.Time) string {
	location, err := time.LoadLocation(DefaultTimeZone)
	if err != nil {
		return value.Format(DateTimeFormat)
	}

	return value.In(location).Format(DateTimeFormat)
}

// FormatDate formats calendar dates consistently.
func FormatDate(value time.Time) string {
	location, err := time.LoadLocation(DefaultTimeZone)
	if err != nil {
		return value.Format(DateFormat)
	}

	return value.In(location).Format(DateFormat)
}

// ParseDate parses a date using the application's standard date format.
func ParseDate(value string) (time.Time, error) {
	location, err := time.LoadLocation(DefaultTimeZone)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation(DateFormat, value, location)
}
