package models

import (
	"fmt"
	"strconv"
	"strings"
)

type CurrencyCode struct {
	Code     string `json:"code"`
	Decimals int64  `json:"decimals"`
}

var allCurrencyCodes []CurrencyCode

func newCurrencyCode(code string, decimals int64) CurrencyCode {
	currencyCode := CurrencyCode{
		Code:     code,
		Decimals: decimals,
	}
	allCurrencyCodes = append(allCurrencyCodes, currencyCode)
	return currencyCode
}

var (
	CurrencyCodeEUR     = newCurrencyCode("EUR", 2)
	CurrencyCodeUSD     = newCurrencyCode("USD", 2)
	CurrencyCodeSEK     = newCurrencyCode("SEK", 2)
	CurrencyCodeDefault = CurrencyCodeSEK
)

func CurrencyFromString(code string) CurrencyCode {
	code = strings.ToUpper(code)
	for _, currency := range allCurrencyCodes {
		if currency.Code == code {
			return currency
		}
	}
	return CurrencyCodeDefault
}

func (c *CurrencyCode) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, c.Code)), nil
}

func (c *CurrencyCode) UnmarshalJSON(data []byte) error {
	code, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	*c = CurrencyFromString(code)
	return nil
}
