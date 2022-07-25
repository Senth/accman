package app

import (
	"github.com/Senth/accman/internal/gateways/sie"
	"strings"
)

func (i *impl) SIEImport(path string) error {
	importer := sie.NewImporter()
	fy, err := importer.Import(path)
	if err != nil {
		return err
	}

	return i.verRepo.AddFiscalYear(*fy)
}

func (i *impl) SIEExport(year string, path string) error {
	exporter := sie.NewExporter(i.accRepo.GetAll())
	fys, err := i.verRepo.GetAll()
	if err != nil {
		return err
	}

	if path == "" {
		path = year + ".se"
	}
	if !strings.HasSuffix(path, ".se") {
		path = path + ".se"
	}

	return exporter.Export(fys, year, path)
}
