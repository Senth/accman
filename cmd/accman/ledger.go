package main

import (
	"fmt"
	"github.com/Senth/accman/internal/gateways/verrepo"
	"github.com/Senth/accman/models"
	"github.com/spf13/cobra"
	"log"
)

var ledgerCmd = &cobra.Command{
	Use:   "ledger year",
	Short: "Print the ledger for the specified year",
	Args:  cobra.ExactArgs(1),
	Run:   ledger,
}

func initLedger(rootCmd *cobra.Command) {
	rootCmd.AddCommand(ledgerCmd)
}

func ledger(cmd *cobra.Command, args []string) {
	repo := verrepo.NewJSONImpl()
	fy := repo.GetFiscalYear(args[0])
	if fy == nil {
		log.Fatalln("Fiscal year not found")
	}

	var balances models.AccountBalances
	balances = append(balances, fy.CalculateEndingBalances()...)
	balances = append(balances, fy.CalculateResult()...)

	for _, b := range balances {
		fmt.Printf("%d: %s\n", b.AccountNumber, b.Balance.Format(models.CurrencyCodeDefault))
	}
}
