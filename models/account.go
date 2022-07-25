package models

import "sort"

type AccountNumber int

type Account struct {
	Number            AccountNumber `json:"number"`
	Name              string        `json:"name"`
	SRUs              []SRU         `json:"srus"`
	VATCode           int           `json:"vatCode"`
	VATPercentage     float64       `json:"vatPercentage"`
	VATAccount        AccountNumber `json:"vatAccount"`
	VATReverseAccount AccountNumber `json:"vatReverseAccount"`
}

type SRU struct {
	Number int `json:"number"`
}

type AccountBalance struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Balance       AmountValue   `json:"balance"`
}

type AccountBalances []AccountBalance

func (a AccountBalances) SortByAccountNumber() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].AccountNumber < a[j].AccountNumber
	})
}
