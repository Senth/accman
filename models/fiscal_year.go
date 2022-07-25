package models

import "sort"

type FiscalYear struct {
	From Date `json:"from"`
	To   Date `json:"to"`
	// Locked fiscal years cannot be changed without unlocking it first
	Locked           bool            `json:"locked"`
	StartingBalances AccountBalances `json:"startingBalances"`
	Verifications    Verifications   `json:"verifications"`
	Changed          bool            `json:"-"`
}

type FiscalYears []FiscalYear

func (f *FiscalYear) AddVerification(verification Verification) {
	f.Verifications = append(f.Verifications, verification)
	f.Changed = true
}

func (f FiscalYear) CalculateEndingBalances() AccountBalances {
	return f.calculate(f.StartingBalances, func(t Transaction) bool {
		return !t.Deleted && t.IsBalance()
	})
}

func (f FiscalYear) CalculateResult() AccountBalances {
	return f.calculate(nil, func(t Transaction) bool {
		return !t.Deleted && t.IsResult()
	})
}

func (f FiscalYear) calculate(input AccountBalances, useTransaction func(t Transaction) bool) AccountBalances {
	balanceMap := make(map[AccountNumber]AmountValue)
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

func (f FiscalYears) SortByDate() {
	sort.Slice(f, func(i, j int) bool {
		return f[i].From.Before(f[j].From)
	})
}

func (f FiscalYears) GetIndex(year string) int {
	for i, y := range f {
		if y.From.Year() == year {
			return i
		}
	}
	return -1
}

// Commit and saves the fiscal year, nothing can be changed without creating new transactions
func (f *FiscalYear) Commit() {
	f.Verifications.SortByDate()
	num := f.Verifications.GetNextNumber()
	for i, v := range f.Verifications {
		if v.Number == 0 {
			v.Commit(num)
			num++
			f.Changed = true
			f.Verifications[i] = v
		}
	}
}
