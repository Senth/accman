package accrepo

import "github.com/Senth/accman/models"

type AccRepo interface {
	Get(accountNumber models.AccountNumber) *models.Account
	GetAll() []models.Account
}
