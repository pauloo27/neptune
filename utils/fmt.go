package utils

import (
	"fmt"
	"strconv"
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

func Pad(n int, minSize int) string {
	str := strconv.Itoa(n)
	for len(str) < minSize {
		str = "0" + str
	}
	return str
}

func FormatDuration(durationInSeconds float64) string {
	minutes := int(durationInSeconds / 60.0)
	seconds := int(durationInSeconds) % 60
	return Fmt("%s:%s", Pad(minutes, 2), Pad(seconds, 2))
}
