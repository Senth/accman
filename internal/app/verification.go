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
	for _, vInfo := range vInfos {
		ver := i.verificationInfoToVerification(vInfo)
		accountFrom := i.getAccount(vInfo.AccountFrom)
		accountTo := i.getAccount(vInfo.AccountTo)
		ver.Transactions = i.createTransactions(accountFrom, accountTo, vInfo.Type, vInfo.Amount)
		err := i.verRepo.AddVerification(ver)
		if err != nil {
			return err
		}
	}
	return nil
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
	// Regular transfer
	if Type.Contains(models.VerificationTypeTransfer) {
		tFrom := models.NewTransaction(accountFrom.Number, amount.Negate())
		tTo := models.NewTransaction(accountTo.Number, amount)
		transactions = append(transactions, tFrom, tTo)
	} else {
		// TODO handle other verification types
		log.Fatalf("currently unsupported verification type: %d", Type)
	}

	return
}
