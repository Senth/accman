package app

import "github.com/Senth/accman/internal/gateways/sie"

func (i *impl) SIEImport(path string) error {
	importer := sie.NewImporter()
	fy, err := importer.Import(path)
	if err != nil {
		return err
	}

	return i.verRepo.AddFiscalYear(*fy)
}
