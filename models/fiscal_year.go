package models

type FiscalYear struct {
	From             string           `json:"from"`
	To               string           `json:"to"`
	StartingBalances []AccountBalance `json:"startingBalances"`
	Verifications    []Verification   `json:"verifications"`
}
