package main

import (
	"fmt"
	"github.com/Senth/accman/internal/gateways/verrepo"
	"github.com/Senth/accman/models"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "accman",
	Short: "Account manager",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		println("\nVerificationTypes")
		printVerificationTypes()
	},
}

func init() {
	initParse(rootCmd)
	initSie(rootCmd)
	initLedger(rootCmd)
	rootCmd.AddCommand(fixCmd)
}

func main() {
	_ = rootCmd.Execute()
}

func printVerificationTypes() {
	fmt.Printf("Invoice in: %d\n", models.VerificationTypeInvoice|models.VerificationTypeIn)
	fmt.Printf("Invoice out: %d\n", models.VerificationTypeInvoice|models.VerificationTypeOut)
	fmt.Printf("Payment in: %d\n", models.VerificationTypePayment|models.VerificationTypeIn)
	fmt.Printf("Payment out: %d\n", models.VerificationTypePayment|models.VerificationTypeOut)
	fmt.Printf("Direct payment in: %d\n", models.VerificationTypeDirectPayment|models.VerificationTypeIn)
	fmt.Printf("Direct payment out: %d\n", models.VerificationTypeDirectPayment|models.VerificationTypeOut)
	fmt.Printf("Transfer: %d\n", models.VerificationTypeTransfer)
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix the database",
	Run: func(cmd *cobra.Command, args []string) {
		repo := verrepo.NewJSONImpl()
		fys, _ := repo.GetAll()
		for _, fy := range fys {
			if fy.From.Year() == "2021" {
				fy = fix(fy)
				_ = repo.UpdateFiscalYear(fy)
			}
		}
	},
}

func fix(fy models.FiscalYear) models.FiscalYear {
	fy.Commit()

	// Set date to the 10th of the month
	for i, v := range fy.Verifications {
		date := v.Date
		date = date.SetDay("10")
		date = date.AddMonths(1)
		v.DateFiled = date
		fy.Verifications[i] = v
	}

	return fy
}
