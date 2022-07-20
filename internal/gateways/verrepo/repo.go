package verrepo

import "github.com/Senth/accman/models"

type VerRepo interface {
	Add(verification ...models.Verification) error
	GetAll() ([]models.FiscalYear, error)
}
