package models

type Transaction struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Amount        Amount        `json:"amount"`
	Created       string        `json:"created"`
	Modified      string        `json:"modified"`
	Deleted       string        `json:"deleted"`
}
