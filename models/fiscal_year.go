package models

type FiscalYear struct {
	From             Date            `json:"from"`
	To               Date            `json:"to"`
	StartingBalances AccountBalances `json:"startingBalances"`
	Verifications    Verifications   `json:"verifications"`
	Changed          bool            `json:"-"`
}

func (f *FiscalYear) AddVerification(verification Verification) {
	f.Verifications = append(f.Verifications, verification)
	f.Changed = true
}

func (f FiscalYear) CalculateEndingBalances() AccountBalances {
	return f.calculate(f.StartingBalances, func(t Transaction) bool {
		return !t.IsDeleted() && t.IsBalance()
	})
}

func (f FiscalYear) CalculateResult() AccountBalances {
	return f.calculate(nil, func(t Transaction) bool {
		return !t.IsDeleted() && t.IsResult()
	})
}

func (f FiscalYear) calculate(input AccountBalances, useTransaction func(t Transaction) bool) AccountBalances {
	balanceMap := make(map[AccountNumber]int64)
	for _, b := range input {
		balanceMap[b.AccountNumber] = b.Balance
	}

	// Sum all the verifications and transactions throughout the year
	for _, v := range f.Verifications {
		for _, t := range v.Transactions {
			// Skip deleted transactions
			if !useTransaction(t) {
				continue
			}

			amount := balanceMap[t.AccountNumber]
			amount += t.Amount.InLocalCurrency()
			balanceMap[t.AccountNumber] = amount
		}
	}

	// Convert the map to a slice, skipping empty accounts
	balances := make(AccountBalances, 0)
	for number, amount := range balanceMap {
		if amount == 0 {
			continue
		}

		balances = append(balances, AccountBalance{
			AccountNumber: number,
			Balance:       amount,
		})
	}

	balances.SortByAccountNumber()
	return balances
}
