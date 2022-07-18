package main

import (
	"fmt"
	"github.com/Senth/accman/models"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "accman",
	Short: "Account manager",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	initParse(rootCmd)
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
