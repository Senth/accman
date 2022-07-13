package models

type Verification struct {
	Name         string           `json:"name"`
	Number       int              `json:"number"`
	Date         string           `json:"date"`
	DateFiled    string           `json:"date_filed"`
	Type         VerificationType `json:"type"`
	Description  string           `json:"description"`
	TotalAmount  Amount           `json:"total_amount"`
	Transactions []Transaction    `json:"transactions"`
}

type VerificationType uint64

const (
	VerificationTypeUnknown VerificationType = 0
	VerificationTypeInvoice VerificationType = 1 << iota
	VerificationTypePayment
	VerificationTypeIn
	VerificationTypeOut
	VerificationTypeTransfer
)
