package sie

import "github.com/Senth/accman/models"

type Importer interface {
	Import(path string) (*models.FiscalYear, error)
}

type Exporter interface {
	Export(fys models.FiscalYears, year string, path string) error
}
