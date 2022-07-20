package parser

import (
	"fmt"
	"github.com/Senth/accman/models"
)

type Parser interface {
	// Verification parses the file into verifications, will return an error if the file doesn't exist
	// or if the file is not valid. But it will also return an error if no parser could be found for the file
	Verification(path string) ([]models.VerificationInfo, error)
	// Text parses the file into a text string
	Text(path string) (Output, error)
}

type Output struct {
	BodyRaw    string
	BodyLayout string
}

type MissingParserError struct {
	Path string
}

func (e MissingParserError) Error() string {
	return fmt.Sprintf("No parser found for file: %s", e.Path)
}
