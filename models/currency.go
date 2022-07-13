package models

type CurrencyCode struct {
	Code     string `json:"code"`
	Decimals int    `json:"decimals"`
}

func newCurrencyCode(code string, decimals int) CurrencyCode {
	return CurrencyCode{
		Code:     code,
		Decimals: decimals,
	}
}

var (
	CurrencyCodeEUR = newCurrencyCode("EUR", 2)
	CurrencyCodeUSD = newCurrencyCode("USD", 2)
	CurrencyCodeSEK = newCurrencyCode("SEK", 2)
)
