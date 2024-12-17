package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCardNumber(t *testing.T) {
	validator := &validator{}

	tests := []struct {
		name            string
		cardNumberInput string
		expectedErr     error
		isCardValid     bool
	}{
		{
			name:            "valid card number",
			cardNumberInput: "5237251624778133",
			expectedErr:     nil,
			isCardValid:     true,
		},
		{
			name:            "valid card number with space",
			cardNumberInput: "52 37251 62477 8133",
			expectedErr:     nil,
			isCardValid:     true,
		},
		{
			name:            "below minimum length",
			cardNumberInput: "12345678901", // 11 digits
			expectedErr:     ErrOutOfRange,
			isCardValid:     false,
		},
		{
			name:            "above maximum length",
			cardNumberInput: "12345678901234567890", // 20 digits
			expectedErr:     ErrOutOfRange,
			isCardValid:     false,
		},
		{
			name:            "card number has invalid characters",
			cardNumberInput: "52h7251624f78133",
			expectedErr:     ErrInvalidCharacters,
			isCardValid:     false,
		},
		{
			name:            "invalid card number (fails Luhn check)",
			cardNumberInput: "5237251624778132",
			expectedErr:     nil,
			isCardValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, err := validator.CardNumberValidator(tt.cardNumberInput)
			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr), "expected error %v, got %v", tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.isCardValid, isValid)
		})
	}
}

func TestValidatorImpl_ValidateCardSchema(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		cardNumber string
		expected   string
		err        error
	}{
		{"American Express", "378282246310005", "American Express", nil},
		{"Visa", "4012888888881881", "Visa", nil},
		{"MasterCard", "5105105105105100", "MasterCard", nil},
		{"JCB", "3530111333300000", "JCB", nil},
		{"Maestro", "6759649826438453", "Maestro", nil},
		{"Unknown scheme", "9999888877776666", "", ErrUnknownScheme},
		{"Invalid characters", "4012X88888881881", "", ErrInvalidCharacters},
		{"below minimum length", "12345678901", "", ErrOutOfRange},
		{"above maximum length", "12345678901234567890", "", ErrOutOfRange},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := validator.CardSchemaValidator(tt.cardNumber)
			assert.Equal(t, tt.expected, schema)
			if tt.err != nil {
				assert.True(t, errors.Is(err, tt.err), "expected error %v, got %v", tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
