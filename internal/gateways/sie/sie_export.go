package sie

import "github.com/Senth/accman/models"

type exporter struct {
}

func NewExporter() Exporter {
	return &exporter{}
}

func (i *exporter) Export(fys []models.FiscalYear, year string, path string) error {
	return nil
}
