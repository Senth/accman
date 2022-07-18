package models

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Amount struct {
	Currency     CurrencyCode `json:"currency"`
	Amount       int64        `json:"amount"`
	LocalAmount  int64        `json:"localAmount"`
	ExchangeRate float64      `json:"exchangeRate"`
}

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
		Amount:   parsedAmount,
	}
}
