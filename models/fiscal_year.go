package models

type FiscalYear struct {
	From             string           `json:"from"`
	To               string           `json:"to"`
	StartingBalances []AccountBalance `json:"starting_balances"`
}
