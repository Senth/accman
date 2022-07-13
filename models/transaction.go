package models

type Transaction struct {
	AccountNumber AccountNumber `json:"account_number"`
	Amount        Amount        `json:"amount"`
	Created       string        `json:"created"`
	Modified      string        `json:"modified"`
	Deleted       string        `json:"deleted"`
}
