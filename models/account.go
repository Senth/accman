package models

type AccountNumber int

type Account struct {
	Number            AccountNumber  `json:"number"`
	Name              string         `json:"name"`
	VATCode           int            `json:"vat_code"`
	VATPercentage     float64        `json:"vat_percentage"`
	VATAccount        *AccountNumber `json:"vat_account"`
	VATReverseAccount *AccountNumber `json:"vat_reverse_account"`
}

type AccountBalance struct {
	AccountNumber AccountNumber `json:"account_number"`
	Balance       int64         `json:"balance"`
}
