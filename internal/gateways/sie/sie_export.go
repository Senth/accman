package sie

import (
	"bytes"
	"fmt"
	"github.com/Senth/accman/models"
	"log"
	"os"
	"time"
)

type exporter struct {
	out      *bytes.Buffer
	accounts []models.Account
	fy       models.FiscalYear
	fys      models.FiscalYears
}

func NewExporter(accounts []models.Account) Exporter {
	return &exporter{
		out:      &bytes.Buffer{},
		accounts: accounts,
	}
}

func (e *exporter) Export(fys models.FiscalYears, year string, path string) error {
	e.addHeader()
	e.addGenerationDate()
	e.addBusinessInfo()
	e.findFiscalYears(fys, year)
	e.addFiscalYears()
	e.addAccounts()
	e.addBalanceAndResults()
	e.addVerifications()
	e.save(path)

	return nil
}

var header = `#FLAGGA 0
#FORMAT PCB
#SIETYP 4
#PROGRAM "accman" 0.1.0
`

func (e *exporter) addHeader() {
	e.out.WriteString(header)
}

func (e *exporter) addGenerationDate() {
	e.out.WriteString(fmt.Sprintf("#GEN %s\n", time.Now().Format(dateFormat)))
}

func (e *exporter) addBusinessInfo() {
	// TODO: add business info
}

func (e *exporter) addAccounts() {
	e.out.WriteString("#KPTYP EUBAS97\n")
	for _, account := range e.accounts {
		// Write for each SRU
		for _, sru := range account.SRUs {
			e.out.WriteString(fmt.Sprintf("#KONTO %d \"%s\"\n", account.Number, account.Name))
			e.out.WriteString(fmt.Sprintf("#KTYP %d %s\n", account.Number, getAccountType(account)))
			e.out.WriteString(fmt.Sprintf("#SRU %d %d\n", account.Number, sru.Number))

			// VAT code
			if account.VATCode != 0 {
				e.out.WriteString(fmt.Sprintf("#MOMSKOD %d %d\n", account.Number, account.VATCode))
			}
		}
	}
}

func getAccountType(account models.Account) string {
	switch {
	case account.Number < 1000:
		log.Fatalf("account number %d is too low", account.Number)
	case account.Number <= 1999:
		return "T"
	case account.Number <= 2999:
		return "S"
	case account.Number <= 3999:
		return "I"
	case account.Number <= 8999:
		return "K"
	}
	log.Fatalf("account number %d is too high", account.Number)
	return ""
}

func (e *exporter) findFiscalYears(fys models.FiscalYears, year string) {
	fys.SortByDate()

	fyIndex := fys.GetIndex(year)
	if fyIndex == -1 {
		log.Fatalf("failed to find fiscal year %s", year)
	}

	e.fy = fys[fyIndex]

	// Get previous three fiscal years to include in the export
	minIndex := fyIndex - 3
	if minIndex < 0 {
		minIndex = 0
	}
	for i := fyIndex; i >= minIndex; i-- {
		e.fys = append(e.fys, fys[i])
	}
}

func (e *exporter) addFiscalYears() {
	for i, fy := range e.fys {
		e.out.WriteString(fmt.Sprintf("#RAR %d %s %s\n", -i, formatDate(fy.From), formatDate(fy.To)))
	}
}

func (e *exporter) addBalanceAndResults() {
	for i, fy := range e.fys {
		e.addBalanceAndResult(fy, -i)
	}
}

func (e *exporter) addBalanceAndResult(fy models.FiscalYear, index int) {
	e.addBalances("#IB", index, fy.StartingBalances)
	e.addBalances("#UB", index, fy.CalculateEndingBalances())
	e.addBalances("#RES", index, fy.CalculateResult())
}

func (e *exporter) addBalances(prefix string, index int, accountBalances models.AccountBalances) {
	for _, accountBalance := range accountBalances {
		number := accountBalance.AccountNumber
		value := accountBalance.Balance.Format(models.CurrencyCodeDefault)
		e.out.WriteString(fmt.Sprintf("%s %d %d %s\n", prefix, index, number, value))
	}
}

func (e *exporter) addVerifications() {
	for _, v := range e.fy.Verifications {
		verString := fmt.Sprintf("#VER A %d %s \"%s\"", v.Number, formatDate(v.Date), v.Name)
		// Add filed date if not same as verification date
		if v.Date != v.DateFiled {
			verString += " " + formatDate(v.DateFiled)
		}
		e.out.WriteString(verString + "\n")
		e.out.WriteString("{\n")

		// Write transactions
		for _, t := range v.Transactions {
			e.addTransaction(v, t)
		}

		e.out.WriteString("}\n")
	}
}

func (e *exporter) addTransaction(v models.Verification, t models.Transaction) {
	if t.Deleted {
		e.out.WriteString(fmt.Sprintf("#BTRANS %d {} %s %s\n", t.AccountNumber, t.Amount.FormatInLocalCurrency(), formatDate(t.Date)))
	} else if isNewTransaction(v, t) {
		e.out.WriteString(fmt.Sprintf("#TRANS %d {} %s\n", t.AccountNumber, t.Amount.FormatInLocalCurrency()))
	} else if isEditTransaction(v, t) {
		e.out.WriteString(fmt.Sprintf("#RTRANS %d {} %s %s\n", t.AccountNumber, t.Amount.FormatInLocalCurrency(), formatDate(t.Date)))
	} else {
		log.Fatalf("unknown transaction, is neither new, edited, or deleted")
	}
}

func isNewTransaction(v models.Verification, t models.Transaction) bool {
	return (t.Date == "" || t.Date == v.DateFiled) && !t.Deleted
}

func isEditTransaction(v models.Verification, t models.Transaction) bool {
	return t.Date != "" && t.Date != v.DateFiled && !t.Deleted
}

const dateFormat = "20060102"

func formatDate(date models.Date) string {
	return date.Time().Format(dateFormat)
}

func (e *exporter) save(path string) {
	if err := os.WriteFile(path, e.out.Bytes(), 0644); err != nil {
		log.Fatalf("failed to write file %s: %v", path, err)
	}
}
