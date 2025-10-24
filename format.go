package main

import (
	"fmt"
	"time"
)

func FormatTaskAge(createdAt time.Time) string {
	duration := time.Since(createdAt)
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)

	if days > 0 {
		return fmt.Sprintf("%dd", days)
	} else if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	} else {
		return "just now"
	}
}
