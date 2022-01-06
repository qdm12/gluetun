package helpers

import "time"

// StringPtr returns a pointer to the string value
// passed as argument.
func StringPtr(s string) *string { return &s }

// DurationPtr returns a pointer to the duration value
// passed as argument.
func DurationPtr(d time.Duration) *time.Duration { return &d }
