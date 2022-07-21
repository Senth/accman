package sie

import (
	"bufio"
	"encoding/json"
	"github.com/Senth/accman/models"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type importer struct {
	fy *models.FiscalYear
	// ver ongoing verification that is being parsed
	ver                 *models.Verification
	endingBalances      models.AccountBalances
	resultBalances      models.AccountBalances
	skipNextTransaction bool
}

func NewImporter() Importer {
	return &importer{}
}

func (i *importer) Import(path string) (*models.FiscalYear, error) {
	i.fy = &models.FiscalYear{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close fiscal year file: %v, err: %v", path, err)
		}
	}(file)

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i.parseLine(line)
	}

	i.verifyBalances()

	return i.fy, nil
}

func (i *importer) parseLine(line string) {
	switch {
	case strings.HasPrefix(line, "#RAR 0"):
		i.setDateRange(line)
	case strings.HasPrefix(line, "#IB 0"):
		i.addStartingBalance(line)
	case strings.HasPrefix(line, "#UB 0"):
		i.addEndingBalance(line)
	case strings.HasPrefix(line, "#RES 0"):
		i.addResultBalance(line)
	case strings.HasPrefix(line, "#VER"):
		i.startVerification(line)
	case strings.HasPrefix(line, "#TRANS"):
		i.addTransaction(line)
	case strings.HasPrefix(line, "#BTRANS"):
		i.addEditedTransaction(line)
	case strings.HasPrefix(line, "#RTRANS"):
		i.addEditedTransaction(line)
		i.skipNextTransaction = true
	case strings.HasPrefix(line, "}"):
		i.endVerification()
	}
}

func (i *importer) setDateRange(line string) {
	re := regexp.MustCompile(`^#RAR 0 (\d{8}) (\d{8})`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		log.Fatalln("invalid date range line: ", line)
	}

	i.fy.From = parseDate(matches[1])
	i.fy.To = parseDate(matches[2])
}

func (i *importer) addStartingBalance(line string) {
	i.fy.StartingBalances = append(i.fy.StartingBalances, parseBalance(line))
}

func (i *importer) addEndingBalance(line string) {
	i.endingBalances = append(i.endingBalances, parseBalance(line))
}

func (i *importer) addResultBalance(line string) {
	i.resultBalances = append(i.resultBalances, parseBalance(line))
}

func parseBalance(line string) models.AccountBalance {
	re := regexp.MustCompile(`^#\w+ 0 (\d{4}) (-?\d+\.\d{2})`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		log.Fatalln("invalid starting balance line: ", line)
	}

	accountNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatalln("invalid account number: ", matches[1])
	}
	amountString := matches[2]
	amount := models.ParseAmount(amountString, models.CurrencyCodeDefault)

	return models.AccountBalance{
		AccountNumber: models.AccountNumber(accountNumber),
		Balance:       amount.Amount,
	}
}

func (i *importer) startVerification(line string) {
	re := regexp.MustCompile(`^#VER A (\d+) (\d{8}) "(.*)" ?(\d{8})?`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		log.Fatalln("invalid verification line: ", line)
	}

	verificationNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatalln("invalid verification number: ", matches[1])
	}

	date := parseDate(matches[2])
	name := matches[3]
	dateFiled := date

	// Filed on a different date
	if matches[4] != "" {
		dateFiled = parseDate(matches[4])
	}

	i.ver = &models.Verification{
		Name:         name,
		Number:       verificationNumber,
		Date:         date,
		DateFiled:    dateFiled,
		Transactions: nil,
	}
}

func (i *importer) addTransaction(line string) {
	if i.skipNextTransaction {
		i.skipNextTransaction = false
		return
	}

	re := regexp.MustCompile(`^#TRANS (\d{4}) {.*?} (-?\d+\.\d{2})`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		log.Fatalln("invalid transaction line: ", line)
	}

	accountNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatalln("invalid account number: ", matches[1])
	}
	amount := models.ParseAmount(matches[2], models.CurrencyCodeDefault)

	transaction := models.Transaction{
		AccountNumber: models.AccountNumber(accountNumber),
		Amount:        amount,
		Created:       i.ver.DateFiled,
	}

	i.ver.Transactions = append(i.ver.Transactions, transaction)
}

func (i *importer) addEditedTransaction(line string) {
	re := regexp.MustCompile(`^#(\w)TRANS (\d{4}) {.*?} (-?\d+\.\d{2}) (\d{8}) "" "" ".*"`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 5 {
		log.Fatalln("invalid edited transaction line: ", line)
	}

	Type := matches[1]
	accountNumber, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatalln("invalid account number: ", matches[2])
	}
	amount := models.ParseAmount(matches[3], models.CurrencyCodeDefault)

	date := parseDate(matches[4])

	createdDate := i.ver.DateFiled
	deletedDate := models.Date("")
	if Type == "B" {
		deletedDate = date
	} else if Type == "R" {
		createdDate = date
	}

	transaction := models.Transaction{
		AccountNumber: models.AccountNumber(accountNumber),
		Amount:        amount,
		Created:       createdDate,
		Deleted:       deletedDate,
	}

	i.ver.Transactions = append(i.ver.Transactions, transaction)
}

func (i *importer) endVerification() {
	err := i.ver.ValidateTransactions()
	if err != nil {
		log.Fatalf("invalid verification %d, err: %v", i.ver.Number, err)
	}

	i.fy.Verifications = append(i.fy.Verifications, *i.ver)
	i.ver = nil
}

func parseDate(dateString string) models.Date {
	date, err := time.Parse("20060102", dateString)
	if err != nil {
		log.Fatalln("invalid date: ", dateString)
	}
	return models.DateFromTime(date)
}

func (i *importer) verifyBalances() {
	endingBalances := i.fy.CalculateEndingBalances()
	i.endingBalances.SortByAccountNumber()
	i.verifyBalance(endingBalances, i.endingBalances)

	resultBalances := i.fy.CalculateResult()
	i.resultBalances.SortByAccountNumber()
	i.verifyBalance(resultBalances, i.resultBalances)
}

func (i *importer) verifyBalance(calculated, imported models.AccountBalances) {
	failed := false
	if len(calculated) != len(imported) {
		failed = true
	}

	for i := range calculated {
		if calculated[i].AccountNumber != imported[i].AccountNumber {
			failed = true
			break
		}
		if calculated[i].Balance != imported[i].Balance {
			failed = true
			break
		}
	}

	// Print diff if failed
	if failed {
		dmp := diffmatchpatch.New()
		calculatedJSON, err := json.Marshal(calculated)
		if err != nil {
			log.Fatalln("failed to marshal calculated balances: ", err)
		}

		importedJSON, err := json.Marshal(imported)
		if err != nil {
			log.Fatalln("failed to marshal imported balances: ", err)
		}

		diffs := dmp.DiffMain(string(calculatedJSON), string(importedJSON), true)
		log.Println(dmp.DiffPrettyText(diffs))
		log.Fatalln("balances do not match")
	}
}
