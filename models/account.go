package models

type AccountNumber int

type Account struct {
	Number            AccountNumber `json:"number"`
	Name              string        `json:"name"`
	VATCode           int           `json:"vatCode"`
	VATPercentage     float64       `json:"vatPercentage"`
	VATAccount        AccountNumber `json:"vatAccount"`
	VATReverseAccount AccountNumber `json:"vatReverseAccount"`
}

type AccountBalance struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Balance       int64         `json:"balance"`
}
