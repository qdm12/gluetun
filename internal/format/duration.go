package format

import (
	"fmt"
	"time"
)

// FriendlyDuration formats a duration in an approximate, human friendly duration.
// For example 55 hours will result in "2 days".
func FriendlyDuration(duration time.Duration) string {
	const twoDays = 48 * time.Hour
	switch {
	case duration < time.Minute:
		seconds := int(duration.Round(time.Second).Seconds())
		const two = 2
		if seconds < two {
			return fmt.Sprintf("%d second", seconds)
		}
		return fmt.Sprintf("%d seconds", seconds)
	case duration <= time.Hour:
		minutes := int(duration.Round(time.Minute).Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	case duration < twoDays:
		hours := int(duration.Truncate(time.Hour).Hours())
		return fmt.Sprintf("%d hours", hours)
	default:
		const hoursInDay = 24
		days := int(duration.Truncate(time.Hour).Hours() / hoursInDay)
		return fmt.Sprintf("%d days", days)
	}
}
