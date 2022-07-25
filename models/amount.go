package models

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Amount struct {
	Currency     CurrencyCode `json:"currency,omitempty"`
	Amount       AmountValue  `json:"amount"`
	LocalAmount  AmountValue  `json:"localAmount,omitempty"`
	ExchangeRate float64      `json:"exchangeRate,omitempty"`
}

type AmountValue int64

var (
	amountParseRegexp = regexp.MustCompile(`^(-?\d{0,3})?[. ,']?(-?\d{0,3})[.,](\d{2})$|^(-?\d{0,3})?[. ,']?(\d{0,3})$`)
	replacement       = `$1$2$4$5.$3`
)

func ParseAmount(amount string, currency CurrencyCode) Amount {
	amount = amountParseRegexp.ReplaceAllString(amount, replacement)
	amountParts := strings.Split(amount, ".")
	var whole, penny string
	if len(amountParts) == 1 {
		whole = amountParts[0]
	}
	if len(amountParts) == 2 {
		whole = amountParts[0]
		penny = amountParts[1]
	}

	// Whole
	if len(whole) == 0 {
		whole = "0"
	}
	parsedAmount, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return Amount{}
	}
	parsedAmount *= int64(math.Pow10(int(currency.Decimals)))

	// Pennies
	if len(penny) == 0 {
		penny = "0"
	}
	parsedPenny, err := strconv.ParseInt(penny, 10, 64)
	if err != nil {
		return Amount{}
	}
	if strings.HasPrefix(amount, "-") {
		parsedPenny *= -1
	}
	parsedAmount += parsedPenny

	return Amount{
		Currency: currency,
		Amount:   AmountValue(parsedAmount),
	}
}

func (a Amount) Negate() Amount {
	return Amount{
		Currency:    a.Currency,
		Amount:      -a.Amount,
		LocalAmount: -a.LocalAmount,
	}
}

// Abs returns the absolute value of the amount.
func (a Amount) Abs() Amount {
	if a.Amount < 0 {
		return a.Negate()
	}
	return a
}

// InLocalCurrency return the amount in the local currency
func (a Amount) InLocalCurrency() AmountValue {
	if a.LocalAmount != 0 {
		return a.LocalAmount
	}
	return a.Amount
}

func (a Amount) FormatInLocalCurrency() string {
	return a.InLocalCurrency().Format(CurrencyCodeDefault)
}

func (a AmountValue) Format(code CurrencyCode) string {
	divider := int64(math.Pow10(int(code.Decimals)))
	whole := int64(a) / divider
	penny := int64(a) % divider
	if penny < 0 {
		penny *= -1
	}
	return fmt.Sprintf("%d.%02d", whole, penny)
}
