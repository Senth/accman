package models

import "testing"

func TestParseAmount(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		expected int64
	}{
		{
			name:     "single digit",
			in:       "1",
			expected: 100,
		},
		{
			name:     "lots of digits",
			in:       "123456789",
			expected: 12345678900,
		},
		{
			name:     "negative",
			in:       "-1",
			expected: -100,
		},
		{
			name:     "negative lots of digits",
			in:       "-123456789",
			expected: -12345678900,
		},
		{
			name:     "decimal",
			in:       "1.23",
			expected: 123,
		},
		{
			name:     "negative decimal",
			in:       "-1.23",
			expected: -123,
		},
		{
			name:     "negative decimal with lots of digits",
			in:       "-123456789.12",
			expected: -12345678912,
		},
		{
			name:     "decimal with comma",
			in:       "1,23",
			expected: 123,
		},
		{
			name:     "number with spaces",
			in:       "123 456",
			expected: 12345600,
		},
		{
			name:     "number with spaces and decimal",
			in:       "123 456,23",
			expected: 12345623,
		},
		{
			name:     "number with spaces and decimal and negative",
			in:       "-123 456,23",
			expected: -12345623,
		},
		{
			name:     "number with ' spacing",
			in:       "123'456",
			expected: 12345600,
		},
		{
			name:     "number with , spacing",
			in:       "123,456",
			expected: 12345600,
		},
		{
			name:     "number with . spacing",
			in:       "123.456",
			expected: 12345600,
		},
		{
			name:     "number with . spacing and decimals with ,",
			in:       "123,456.23",
			expected: 12345623,
		},
		{
			name:     "number with . spacing and decimals with .",
			in:       "123.456.23",
			expected: 12345623,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ParseAmount(tc.in, CurrencyCodeSEK)
			if actual.Amount != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual.Amount)
			}
		})
	}
}
