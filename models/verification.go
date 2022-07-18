package models

type Verification struct {
	Name         string           `json:"name"`
	Number       int              `json:"number"`
	Date         string           `json:"date"`
	DateFiled    string           `json:"dateFiled"`
	Type         VerificationType `json:"type"`
	Description  string           `json:"description"`
	TotalAmount  Amount           `json:"totalAmount"`
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
