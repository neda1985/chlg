package validator

import (
	"fmt"
	"strings"
)

type Validator interface {
	CardNumberValidator(cardNumber string) (bool, error)
	CardSchemaValidator(cardNumber string) (string, error)
}
type validator struct{}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) CardNumberValidator(cardNumber string) (bool, error) {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	// I do this general length validation here since we don't want loop through extremely long string.
	if len(cardNumber) < 12 || len(cardNumber) > 19 {
		return false, fmt.Errorf("%w: got %d characters", ErrOutOfRange, len(cardNumber))
	}
	if !isAllDigits(cardNumber) {
		return false, ErrInvalidCharacters
	}
	return luhnCheck(cardNumber), nil
}

func (v *validator) CardSchemaValidator(cardNumber string) (string, error) {
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	if len(cardNumber) < 12 || len(cardNumber) > 19 {
		return "", fmt.Errorf("%w: got %d characters", ErrOutOfRange, len(cardNumber))
	}
	if !isAllDigits(cardNumber) {
		return "", ErrInvalidCharacters
	}
	for _, scheme := range getCardSchemes() {
		matches, err := isValidForScheme(cardNumber, scheme)
		if err != nil {
			return "", fmt.Errorf("error validating scheme %q: %w", scheme.Name, err)
		}
		if matches {
			return scheme.Name, nil
		}
	}

	return "", ErrUnknownScheme
}
