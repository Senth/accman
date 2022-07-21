package models

type Transaction struct {
	AccountNumber AccountNumber `json:"accountNumber"`
	Amount        Amount        `json:"amount"`
	Created       Date          `json:"created,omitempty"`
	Deleted       Date          `json:"deleted,omitempty"`
}

// NewTransaction create a new transaction and sets the created and modified date to now
func NewTransaction(accountNumber AccountNumber, amount Amount) Transaction {
	return Transaction{
		AccountNumber: accountNumber,
		Amount:        amount,
		Created:       DateNow(),
	}
}

func (t Transaction) IsDeleted() bool {
	return t.Deleted != ""
}

func (t Transaction) IsBalance() bool {
	return t.AccountNumber < 3000
}

func (t Transaction) IsResult() bool {
	return t.AccountNumber >= 3000
}
