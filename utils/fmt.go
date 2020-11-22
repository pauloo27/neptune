package utils

import "fmt"

func Fmt(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}

func EnforceSize(text string, maxLen int) string {
	if maxLen == 0 || len(text) <= maxLen {
		return text
	}

	return text[0:maxLen-3] + "..."
}
