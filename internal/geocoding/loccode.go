package geocoding

import (
	"fmt"
	"strings"
)

const locCodeSeparator = "|"

type LocCode struct {
	Continent   string
	CountryCode string
	StateCode   string
	CityName    string
}

func ParseLocCode(raw string) (LocCode, error) {
	parts := extractParts(raw)

	partsLen := len(parts)

	if partsLen != 4 {
		return LocCode{}, fmt.Errorf("invalid locCode: expected 4 parts, got %d", partsLen)
	}

	result := LocCode{
		Continent:   parts[0],
		CountryCode: parts[1],
		StateCode:   parts[2],
		CityName:    parts[3],
	}

	if err := result.validate(); err != nil {
		return LocCode{}, err
	}

	return result, nil
}

func extractParts(raw string) []string {
	trimmedRaw := strings.TrimSpace(raw)
	parts := strings.Split(trimmedRaw, locCodeSeparator)

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}

func (l LocCode) validate() error {
	if l.Continent == "" {
		return fmt.Errorf("continent is required")
	}

	if l.CountryCode == "" {
		return fmt.Errorf("countryCode is required")
	}

	if l.StateCode == "" {
		return fmt.Errorf("stateCode is required")
	}

	if l.CityName == "" {
		return fmt.Errorf("cityName is required")
	}

	return nil
}
