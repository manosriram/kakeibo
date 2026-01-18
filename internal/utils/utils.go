package utils

import (
	"fmt"
	"time"
)

func getDaySuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func FormatDateTime(t time.Time) string {
	t = t.Local()
	day := t.Day()
	suffix := getDaySuffix(day)

	// Format hour for 12-hour clock
	hour := t.Hour()
	minute := t.Minute()
	ampm := "am"

	if hour >= 12 {
		ampm = "pm"
		if hour > 12 {
			hour -= 12
		}
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%d%s %s, %d %d:%02d%s",
		day, suffix, t.Month().String(), t.Year(), hour, minute, ampm)
}
