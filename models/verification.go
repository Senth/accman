package models

import (
	"fmt"
	"sort"
)

type Verification struct {
	Name         string           `json:"name"`
	Number       int              `json:"number,omitempty"`
	Date         Date             `json:"date"`
	DateFiled    Date             `json:"dateFiled,omitempty"`
	Type         VerificationType `json:"type"`
	Description  string           `json:"description,omitempty"`
	Transactions []Transaction    `json:"transactions"`
}

type VerificationType uint64

const (
	VerificationTypeInvoice VerificationType = 1 << iota
	VerificationTypePayment
	VerificationTypeDirectPayment
	VerificationTypeIn
	VerificationTypeOut
	VerificationTypeTransfer
	VerificationTypeUnknown VerificationType = 0
)

// Contains return true if all the types is part of this type
func (t VerificationType) Contains(other ...VerificationType) bool {
	for _, o := range other {
		if t&o == 0 {
			return false
		}
	}
	return true
}

// VerificationInfo shorthand information for a verification
// Used when parsing verifications
type VerificationInfo struct {
	Date        Date
	Name        string
	Type        VerificationType
	AccountFrom AccountNumber
	AccountTo   AccountNumber
	Amount      Amount
}

type Verifications []Verification

func (v Verification) ValidateTransactions() error {
	sum := AmountValue(0)
	for _, t := range v.Transactions {
		if !t.Deleted {
			sum += t.Amount.InLocalCurrency()
		}
	}

	if sum != 0 {
		return fmt.Errorf("verification %d has a sum of %d", v.Number, sum)
	}
	return nil
}

func (v *Verification) Commit(number int) {
	v.Number = number
	v.DateFiled = DateNow()
}

func (v Verifications) SortByDate() {
	sort.Slice(v, func(i, j int) bool {
		return v[i].Date.Before(v[j].Date)
	})
}

func (v Verifications) GetNextNumber() int {
	max := 0
	for _, v := range v {
		if v.Number > max {
			max = v.Number
		}
	}
	return max + 1
}
