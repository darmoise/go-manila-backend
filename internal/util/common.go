package util

import (
	"math"
	"strings"
	"time"
)

const TemperatureUnit = "c"
const KmPh = "km/h"

func RoundToInt(value float64) int {
	return int(math.Round(value))
}

func DayName(date string) string {
	t, err := time.Parse(time.DateOnly, date)

	if err != nil {
		return ""
	}

	return t.Format("Mon")
}

func FirstOrZero(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}

	return arr[0]
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}
