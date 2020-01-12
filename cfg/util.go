package cfg

import (
	"time"
)

// Defaults ...
const (
	DateFormat = "2006-01-02"

	Day = time.Hour * 24
)

// FormatDate ...
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// ParseDate ...
func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}

// FormatTime ...
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StartOfDay ...
func StartOfDay(t time.Time) time.Time {
	return t.Truncate(Day)
}

// EndOfDay ...
func EndOfDay(t time.Time) time.Time {
	return t.Round(Day)
}
