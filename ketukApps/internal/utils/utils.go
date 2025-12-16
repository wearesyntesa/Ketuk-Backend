package utils

import (
	"time"
)

func ParseTimeLayout(layout, value string) (time.Time, error) {
	// Implementation of time parsing utility
	return time.Parse(layout, value)
}
