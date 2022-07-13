package models

type Amount struct {
	Code         CurrencyCode `json:"code"`
	Amount       int64        `json:"amount"`
	LocalAmount  int64        `json:"local_amount"`
	ExchangeRate float64      `json:"exchange_rate"`
}
