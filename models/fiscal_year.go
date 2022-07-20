package models

type FiscalYear struct {
	From             Date             `json:"from"`
	To               Date             `json:"to"`
	StartingBalances []AccountBalance `json:"startingBalances"`
	Verifications    Verifications    `json:"verifications"`
	Changed          bool             `json:"-"`
}

func (f *FiscalYear) AddVerification(verification Verification) {
	f.Verifications = append(f.Verifications, verification)
	f.Changed = true
}
