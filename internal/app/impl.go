package app

import (
	"github.com/Senth/accman/internal/gateways/accrepo"
	"github.com/Senth/accman/internal/gateways/parser"
	"github.com/Senth/accman/internal/gateways/verrepo"
)

type impl struct {
	parser  parser.Parser
	accRepo accrepo.AccRepo
	verRepo verrepo.VerRepo
}

func NewApp() App {
	return &impl{
		parser:  parser.NewParser(),
		accRepo: accrepo.NewJSONImpl(),
		verRepo: verrepo.NewJSONImpl(),
	}
}
