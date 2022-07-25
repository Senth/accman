package app

import (
	"github.com/Senth/accman/models"
	"log"
)

func (i impl) VerificationParse(path string) error {
	vInfos, err := i.parser.Verification(path)
	if err != nil {
		return err
	}

	return i.VerificationAdd(vInfos)
}

func (i impl) VerificationAdd(vInfos []models.VerificationInfo) error {
	// Create verifications from information
	verifications := make([]models.Verification, 0, len(vInfos))
	for _, vInfo := range vInfos {
		ver := i.verificationInfoToVerification(vInfo)
		accountFrom := i.getAccount(vInfo.AccountFrom)
		accountTo := i.getAccount(vInfo.AccountTo)
		ver.Transactions = i.createTransactions(accountFrom, accountTo, vInfo.Type, vInfo.Amount)
		verifications = append(verifications, ver)
	}

	// Add verifications
	return i.verRepo.AddVerification(verifications...)
}

func (i impl) verificationInfoToVerification(vInfo models.VerificationInfo) models.Verification {
	return models.Verification{
		Name:      vInfo.Name,
		Date:      vInfo.Date,
		DateFiled: models.DateNow(),
		Type:      vInfo.Type,
	}
}

func (i impl) getAccount(number models.AccountNumber) models.Account {
	account := i.accRepo.Get(number)
	if account == nil {
		log.Fatalf("account number %d not found", number)
	}
	return *account
}

func (i impl) createTransactions(accountFrom, accountTo models.Account, Type models.VerificationType, amount models.Amount) (transactions []models.Transaction) {
	switch {
	// Regular transaction
	case isSimpleTransfer(Type):
		tFrom := models.NewTransaction(accountFrom.Number, amount.Negate())
		tTo := models.NewTransaction(accountTo.Number, amount)
		transactions = append(transactions, tFrom, tTo)
	// Direct payment out
	case Type.Contains(models.VerificationTypeDirectPayment, models.VerificationTypeOut):
		transactions = append(transactions, i.createAdvancedTransaction(accountFrom, accountTo, amount)...)
	// Invoice In
	case Type.Contains(models.VerificationTypeInvoice, models.VerificationTypeIn):
		transactions = append(transactions, i.createAdvancedTransaction(accountFrom, accountTo, amount)...)
	default:
		// TODO handle other verification types
		log.Fatalf("currently unsupported verification type: %d", Type)
	}
	return
}

func isSimpleTransfer(Type models.VerificationType) bool {
	switch {
	case Type.Contains(models.VerificationTypeTransfer):
		return true
	case Type.Contains(models.VerificationTypePayment):
		return true
	default:
		return false
	}
}

func (i impl) createAdvancedTransaction(accountFrom, accountTo models.Account, amount models.Amount) (transactions []models.Transaction) {
	if accountFrom.VATPercentage != 0 || accountTo.VATPercentage != 0 {
		log.Fatalf("VAT handling is not supported yet")
	}

	tFrom := models.NewTransaction(accountFrom.Number, amount.Negate())
	transactions = append(transactions, tFrom)
	tTo := models.NewTransaction(accountTo.Number, amount)
	transactions = append(transactions, tTo)

	return
}
