package sie

import "github.com/Senth/accman/models"

type Importer interface {
	Import(path string) (*models.FiscalYear, error)
}

type Exporter interface {
	Export(all []models.FiscalYear, year string, path string) error
}
