package parser

import (
	"encoding/json"
	"github.com/Senth/accman/models"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const filename = "parser.json"

type Impl struct {
	parsers []parser
}

func NewParser() Parser {
	p := &Impl{}
	p.load()
	return p
}

// load the parsers from the JSON file
func (i *Impl) load() {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicf("failed to open parser file: %v", err)
	}

	err = json.NewDecoder(file).Decode(&i.parsers)
	if err != nil {
		log.Panicf("failed to decode parser file: %v", err)
	}
}

func (i Impl) Verification(path string) ([]models.Verification, error) {
	output, err := i.Text(path)
	if err != nil {
		return nil, err
	}

	parser := i.getParser(output.BodyRaw)
	if parser == nil {
		return nil, &MissingParserError{}
	}

	return parser.parse(output)
}

func (i Impl) getParser(text string) *parser {
	for _, parser := range i.parsers {
		if strings.Contains(text, parser.Identifier) {
			return &parser
		}
	}
	return nil
}

func (i Impl) Text(path string) (Output, error) {
	rawChan, rawErrChan := parsePdf(path, false)
	layoutChan, layoutErrChan := parsePdf(path, true)

	return Output{
		BodyRaw:    <-rawChan,
		BodyLayout: <-layoutChan,
	}, consolidateErrors(rawErrChan, layoutErrChan)
}

func parsePdf(path string, layout bool) (chan string, chan error) {
	bodyChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		args := []string{"-q", "-nopgbrk", "-enc", "UTF-8", "-eol", "unix", path}
		if layout {
			args = append(args, "-layout")
		}
		args = append(args, "-")
		body, err := exec.Command("pdftotext", args...).Output()
		if err != nil {
			errChan <- err
			bodyChan <- ""
			return
		}

		bodyChan <- string(body)
		errChan <- nil
		return
	}()

	return bodyChan, errChan
}

func consolidateErrors(errChannels ...chan error) (err error) {
	for _, errChan := range errChannels {
		err = <-errChan
	}
	return err
}

type parser struct {
	Identifier string  `json:"identifier"`
	Options    options `json:"options"`
	// Regexp containing date, name, and currency groups
	Currency            string              `json:"currency"`
	Regexp              string              `json:"regexp"`
	DateFormat          string              `json:"dateFormat"`
	VerificationParsers verificationParsers `json:"verificationParsers"`
}

type verificationParsers []verificationParser

type options struct {
	PDFLayout pdfLayout `json:"pdfLayout"`
}

type pdfLayout string

const (
	PDFLayoutRaw    pdfLayout = "raw"
	PDFLayoutLayout pdfLayout = "layout"
)

type verificationParser struct {
	Identifier  string                  `json:"identifier"`
	Name        string                  `json:"name"`
	AccountFrom int                     `json:"accountFrom"`
	AccountTo   int                     `json:"accountTo"`
	Type        models.VerificationType `json:"type"`
}

type verificationInfo struct {
	Date        string
	Name        string
	Type        models.VerificationType
	AccountFrom int
	AccountTo   int
	Amount      models.Amount
}

type regexMatch struct {
	Date   string        `json:"date"`
	Name   string        `json:"name"`
	Amount models.Amount `json:"amount"`
}

func (p parser) parse(output Output) ([]models.Verification, error) {

	text := p.getTextToUse(output)
	matches, err := p.getMatches(text)
	if err != nil {
		return nil, err
	}

	var vInfos []verificationInfo
	for _, match := range matches {
		vParser := p.VerificationParsers.find(match.Name)
		if vParser == nil {
			println("Skipping verification since no identifier could be matched for: " + match.Name)
		}
		name := match.Name
		if vParser.Name != "" {
			name = vParser.Name
		}
		verification := verificationInfo{
			Date:        match.Date,
			Name:        name,
			Amount:      match.Amount,
			Type:        vParser.Type,
			AccountFrom: vParser.AccountFrom,
			AccountTo:   vParser.AccountTo,
		}
		vInfos = append(vInfos, verification)
	}

	return nil, nil
}

func (p parser) getTextToUse(output Output) string {
	switch p.Options.PDFLayout {
	case PDFLayoutRaw:
		return output.BodyRaw
	case PDFLayoutLayout:
		return output.BodyLayout
	default:
		return output.BodyLayout
	}
}

func (p parser) getMatches(text string) ([]regexMatch, error) {
	re, err := regexp.Compile(p.Regexp)
	if err != nil {
		return nil, err
	}
	matches := re.FindAllStringSubmatch(text, -1)

	mappedMatches := make([]regexMatch, len(matches))
	for i, match := range matches {
		reMatch := regexMatch{}
		for j, groupName := range re.SubexpNames() {
			value := strings.TrimSpace(match[j])

			switch groupName {
			case "date":
				// Fix date
				if p.DateFormat != "" {
					date, err := time.Parse(p.DateFormat, value)
					if err != nil {
						return nil, err
					}
					value = date.Format("2006-01-02")
				}
				reMatch.Date = value
			case "name":
				reMatch.Name = value
			case "amount":
				// Fix amount
				currency := models.CurrencyFromString(p.Currency)
				reMatch.Amount = models.ParseAmount(value, currency)
			}
		}
		mappedMatches[i] = reMatch
	}
	return mappedMatches, nil
}

func (v verificationParsers) find(name string) *verificationParser {
	for _, parser := range v {
		if parser.Identifier == name {
			return &parser
		}
	}
	return nil
}

func (i verificationInfo) Model() models.Verification {
	return models.Verification{
		Date: i.Date,
		Name: i.Name,
		Type: i.Type,
	}
}
