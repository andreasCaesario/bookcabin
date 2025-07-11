package utils

import (
	"strings"
	"time"
)

// ValidateDateNotPast checks if a date string in DD - MM - YY format is today or in the future.
func ValidateDateNotPast(dateStr string) bool {
	dateParts := strings.Split(dateStr, "-")
	if len(dateParts) != 3 {
		return false
	}
	dd := strings.TrimSpace(dateParts[0])
	mm := strings.TrimSpace(dateParts[1])
	yy := strings.TrimSpace(dateParts[2])
	if len(dd) != 2 || len(mm) != 2 || len(yy) != 2 {
		return false
	}
	fullYear := "20" + yy
	if yy < "50" {
		fullYear = "20" + yy
	} else {
		fullYear = "19" + yy
	}
	inputDateStr := fullYear + "-" + mm + "-" + dd
	inputDate, err := time.Parse("2006-01-02", inputDateStr)
	if err != nil {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return !inputDate.Before(today)
}

// MakeRange returns a slice of ints from start to end inclusive
func MakeRange(start, end int) []int {
	rangeSlice := make([]int, end-start+1)
	for i := range rangeSlice {
		rangeSlice[i] = start + i
	}
	return rangeSlice
}
