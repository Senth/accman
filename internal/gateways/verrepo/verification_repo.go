package verrepo

import "github.com/Senth/accman/models"

type VerRepo interface {
	AddVerification(verification ...models.Verification) error
	AddFiscalYear(fiscalYear models.FiscalYear) error
	GetAll() (models.FiscalYears, error)
	GetFiscalYear(year string) *models.FiscalYear
	UpdateFiscalYear(fiscalYear models.FiscalYear) error
}
