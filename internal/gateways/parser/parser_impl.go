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

type impl struct {
	parsers []parser
}

func NewParser() Parser {
	p := &impl{}
	p.load()
	return p
}

// load the parsers from the JSON file
func (i *impl) load() {
	file, err := os.Open(filename)
	if err != nil {
		log.Panicf("failed to open parser file: %v", err)
	}

	err = json.NewDecoder(file).Decode(&i.parsers)
	if err != nil {
		log.Panicf("failed to decode parser file: %v", err)
	}
}

func (i impl) Verification(path string) (vInfos []models.VerificationInfo, err error) {
	output, err := i.Text(path)
	if err != nil {
		return nil, err
	}

	parsers := i.getParsers(output.BodyRaw)
	if len(parsers) == 0 {
		return nil, &MissingParserError{}
	}

	for _, p := range parsers {
		infos, err := p.parse(output)
		if err != nil {
			return nil, err
		}
		vInfos = append(vInfos, infos...)
	}
	return
}

func (i impl) getParsers(text string) (parsers []parser) {
	for _, parser := range i.parsers {
		if strings.Contains(text, parser.Identifier) {
			parsers = append(parsers, parser)
		}
	}
	return
}

func (i impl) Text(path string) (Output, error) {
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
	// Prefix optional prefix that will be added to all verification names
	Prefix string `json:"prefix"`
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
	Identifier    string                  `json:"identifier"`
	Name          string                  `json:"name"`
	AccountFrom   models.AccountNumber    `json:"accountFrom"`
	AccountTo     models.AccountNumber    `json:"accountTo"`
	Type          models.VerificationType `json:"type"`
	Bidirectional bool                    `json:"bidirectional"`
}

type regexMatch struct {
	Date   string        `json:"date"`
	Name   string        `json:"name"`
	Amount models.Amount `json:"amount"`
}

func (p parser) parse(output Output) ([]models.VerificationInfo, error) {

	text := p.getTextToUse(output)
	matches, err := p.getMatches(text)
	if err != nil {
		return nil, err
	}

	var vInfos []models.VerificationInfo
	for _, match := range matches {
		vParser := p.VerificationParsers.find(match.Name)
		if vParser == nil {
			println("Skipping verification since no identifier could be matched for: " + match.Name)
			continue
		}

		verification := models.VerificationInfo{
			Date:        models.Date(match.Date),
			Name:        p.getVerificationName(*vParser, match.Name),
			Amount:      match.Amount,
			Type:        vParser.Type,
			AccountFrom: vParser.AccountFrom,
			AccountTo:   vParser.AccountTo,
		}

		// Fix bidirectional verifications
		if vParser.Bidirectional && verification.Amount.Amount < 0 {
			verification.AccountTo, verification.AccountFrom = verification.AccountFrom, verification.AccountTo
		}
		verification.Amount = verification.Amount.Abs()

		vInfos = append(vInfos, verification)
	}

	return vInfos, nil
}

func (p parser) getVerificationName(vParser verificationParser, name string) string {
	if vParser.Name != "" {
		name = vParser.Name
	}
	if p.Prefix != "" {
		name = p.Prefix + name
	}
	return name
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
		if strings.Contains(name, parser.Identifier) {
			return &parser
		}
	}
	return nil
}
