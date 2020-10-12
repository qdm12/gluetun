package logging

import (
	"fmt"
	"time"
)

func FormatDuration(duration time.Duration) string {
	switch {
	case duration < time.Minute:
		seconds := int(duration.Round(time.Second).Seconds())
		if seconds < 2 {
			return fmt.Sprintf("%d second", seconds)
		}
		return fmt.Sprintf("%d seconds", seconds)
	case duration <= time.Hour:
		minutes := int(duration.Round(time.Minute).Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	case duration < 48*time.Hour:
		hours := int(duration.Truncate(time.Hour).Hours())
		return fmt.Sprintf("%d hours", hours)
	default:
		days := int(duration.Truncate(time.Hour).Hours() / 24)
		return fmt.Sprintf("%d days", days)
	}
}
