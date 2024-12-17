package validator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrInvalidCharacters = errors.New("card number must contain only digits")
	ErrOutOfRange        = errors.New("card number out of range")
	ErrUnknownScheme     = errors.New("unknown card scheme")
	cardSchemes          []CardScheme
	cardSchemesOnce      sync.Once
)

type CardScheme struct {
	Name           string
	PrefixRanges   []string
	AllowedLengths []int
}

// initializeCardSchemes initializes the list of supported card schemes
func initializeCardSchemes() {
	cardSchemes = []CardScheme{
		{
			Name:           "American Express",
			PrefixRanges:   []string{"34", "37"},
			AllowedLengths: []int{15},
		},
		{
			Name:           "JCB",
			PrefixRanges:   []string{"3528-3589"},
			AllowedLengths: []int{16, 17, 18, 19},
		},
		{
			Name:           "Maestro",
			PrefixRanges:   []string{"50", "56-58", "6"},
			AllowedLengths: []int{12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			Name:           "Visa",
			PrefixRanges:   []string{"4"},
			AllowedLengths: []int{13, 16, 19},
		},
		{
			Name:           "MasterCard",
			PrefixRanges:   []string{"2221-2720", "51-55"},
			AllowedLengths: []int{16},
		},
	}
}

// getCardSchemes ensures cardSchemes is initialized only once and returns it
func getCardSchemes() []CardScheme {
	cardSchemesOnce.Do(initializeCardSchemes)
	return cardSchemes
}

func isAllDigits(cardNumber string) bool {
	for _, char := range cardNumber {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func luhnCheck(cardNumber string) bool {
	sum := 0
	alternate := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// isValidForScheme checks if a card number matches a specific card scheme
func isValidForScheme(cardNumber string, scheme CardScheme) (bool, error) {
	length := len(cardNumber)

	// Check if length is allowed for the scheme
	if !contains(scheme.AllowedLengths, length) {
		return false, nil
	}

	// Check if the prefix matches any of the scheme's prefix ranges
	for _, rangeStr := range scheme.PrefixRanges {
		matches, err := matchesPrefixRange(cardNumber, rangeStr)
		if err != nil {
			return false, fmt.Errorf("error validating prefix range %q: %w", rangeStr, err)
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

// matchesPrefixRange checks if a card number's prefix matches a specific range
func matchesPrefixRange(cardNumber, rangeStr string) (bool, error) {
	if strings.Contains(rangeStr, "-") {
		parts := strings.Split(rangeStr, "-")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid range format: %s", rangeStr)
		}

		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("invalid range values: %s", rangeStr)
		}

		prefixLength := len(parts[0])
		if len(cardNumber) < prefixLength {
			return false, nil
		}

		prefix, err := strconv.Atoi(cardNumber[:prefixLength])
		if err != nil {
			return false, fmt.Errorf("invalid prefix conversion: %w", err)
		}

		return prefix >= start && prefix <= end, nil
	}
	if strings.HasPrefix(cardNumber, rangeStr) {
		return true, nil
	}

	return false, nil
}

func contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
