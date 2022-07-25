package app

import "github.com/Senth/accman/models"

type App interface {
	VerificationParse(path string) error
	VerificationAdd(verificationInfos []models.VerificationInfo) error
	SIEImport(path string) error
	SIEExport(year string, path string) error
}
