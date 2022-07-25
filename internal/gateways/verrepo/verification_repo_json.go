package verrepo

import (
	"encoding/json"
	"fmt"
	"github.com/Senth/accman/models"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type impl struct {
	fys []*models.FiscalYear
}

func NewJSONImpl() VerRepo {
	i := &impl{}
	i.load()
	return i
}

func (i *impl) AddVerification(verification ...models.Verification) error {
	for _, v := range verification {
		fy := i.getFiscalYear(v.Date)
		if fy == nil {
			log.Fatalf("failed to find fiscal year for date: %v", v.Date)
		}

		if fy.Locked {
			log.Println("fiscal year is locked, skipping")
			continue
		}

		fy.AddVerification(v)
	}

	i.save()

	return nil
}

func (i impl) GetAll() (fys models.FiscalYears, err error) {
	fys = make(models.FiscalYears, len(i.fys))
	for i, fy := range i.fys {
		fys[i] = *fy
	}
	return fys, nil
}

func (i *impl) load() {
	files := getFiscalYearFiles()

	for _, file := range files {
		i.loadFiscalYears(file)
	}
}

func (i *impl) sort() {
	for _, fy := range i.fys {
		if fy.Changed {
			fy.Verifications.SortByDate()
		}
	}
}

func (i *impl) save() {
	i.sort()

	for _, fy := range i.fys {
		if fy.Changed {
			backupFiscalYear(fy)
			saveFiscalYear(fy)
		}
	}
}

func getFiscalYearFiles() (files []string) {
	fiscalYearRegex := regexp.MustCompile(`^\d{4}\.json$`)
	// Iterate over the directory and get all the files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Error walking the path: %v", err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-matching files
		if !fiscalYearRegex.MatchString(info.Name()) {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		log.Fatalf("failed to get json fiscal years: %v", err)
	}

	return
}

func saveFiscalYear(fy *models.FiscalYear) {
	file, err := os.Create(getFiscalYearFilename(fy))
	if err != nil {
		log.Fatalf("failed to create fiscal year file: %v, err: %v", getFiscalYearFilename(fy), err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close fiscal year file: %v, err: %v", getFiscalYearFilename(fy), err)
		}
	}(file)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(fy)
	if err != nil {
		log.Fatalf("failed to encode fiscal year file: %v, err: %v", getFiscalYearFilename(fy), err)
	}
}

func (i *impl) loadFiscalYears(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open fiscal year file: %v, err: %v", filepath, err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close fiscal year file: %v, err: %v", filepath, err)
		}
	}(file)

	var fy models.FiscalYear
	err = json.NewDecoder(file).Decode(&fy)
	if err != nil {
		log.Fatalf("failed to decode fiscal year file: %v, err: %v", filepath, err)
	}
	i.fys = append(i.fys, &fy)
}

func (i impl) GetFiscalYear(year string) *models.FiscalYear {
	date := models.Date(year + "-01-01")
	return i.getFiscalYear(date)
}

func (i impl) getFiscalYear(date models.Date) *models.FiscalYear {
	for _, fy := range i.fys {
		if date.Between(fy.From, fy.To) {
			return fy
		}
	}
	return nil
}

func getFiscalYearFilename(fy *models.FiscalYear) string {
	return fmt.Sprintf("%v.json", fy.From.Year())
}

func getFiscalYearBackupName(fy *models.FiscalYear) string {
	timestamp := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", ".")
	filename := fmt.Sprintf("%v %v.json", fy.From.Year(), timestamp)
	return filepath.Join("backup", filename)
}

func backupFiscalYear(fy *models.FiscalYear) {
	// Only backup if a previous file exists
	if _, err := os.Stat(getFiscalYearFilename(fy)); err == nil {
		createBackupDir()
		err := os.Rename(getFiscalYearFilename(fy), getFiscalYearBackupName(fy))
		if err != nil {
			log.Fatalf("failed to backup fiscal year: %v", err)
		}
	}
}

func createBackupDir() {
	err := os.MkdirAll("backup", 0755)
	if err != nil {
		log.Fatalf("failed to create backup directory: %v", err)
	}
}

func (i *impl) AddFiscalYear(fy models.FiscalYear) error {
	for _, f := range i.fys {
		if f.From.Year() == fy.From.Year() {
			return fmt.Errorf("fiscal year already exists: %v", fy.From.Year())
		}
	}

	fy.Changed = true
	i.fys = append(i.fys, &fy)
	i.save()
	return nil
}

func (i *impl) UpdateFiscalYear(fy models.FiscalYear) error {
	for j, f := range i.fys {
		if f.From.Year() == fy.From.Year() {
			i.fys[j] = &fy
			saveFiscalYear(&fy)
			return nil
		}
	}

	return fmt.Errorf("fiscal year does not exist: %v", fy.From.Year())
}
