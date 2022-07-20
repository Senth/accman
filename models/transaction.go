package models

type Transaction struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Amount        Amount        `json:"amount"`
	Created       Date          `json:"created,omitempty"`
	Modified      Date          `json:"modified,omitempty"`
	Deleted       Date          `json:"deleted,omitempty"`
}

// NewTransaction create a new transaction and sets the created and modified date to now
func NewTransaction(accountNumber AccountNumber, amount Amount) Transaction {
	return Transaction{
		AccountNumber: accountNumber,
		Amount:        amount,
		Created:       DateNow(),
		Modified:      DateNow(),
	}
}
