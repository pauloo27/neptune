package utils

import (
	"fmt"
	"time"
)

func Fmt(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}

func EnforceSize(text string, maxLen int) string {
	if maxLen == 0 || len(text) <= maxLen {
		return text
	}

	return text[0:maxLen-3] + "..."
}

func FormatDuration(duration time.Duration) string {
	minutes := int(duration.Seconds() / 60.0)
	seconds := int(duration.Seconds()) % 60
	return Fmt("%d:%d", minutes, seconds)
}
