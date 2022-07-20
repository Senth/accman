package main

import (
	"github.com/Senth/accman/internal/app"
	"github.com/Senth/accman/internal/gateways/parser"
	"github.com/spf13/cobra"
	"log"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a file and either import it to the database or print the parsed text into stdout",
	Run:   parse,
}

func initParse(rootCmd *cobra.Command) {
	rootCmd.AddCommand(parseCmd)
	parseCmd.AddCommand(parsePrintCmd)
	parseCmd.AddCommand(parseImportCmd)
}

func parse(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

var parsePrintCmd = &cobra.Command{
	Use:   "print [files]",
	Short: "Print the parsed text from a file",
	Args:  cobra.MinimumNArgs(1),
	Run:   parsePrint,
}

func parsePrint(cmd *cobra.Command, args []string) {
	p := parser.NewParser()

	for _, path := range args {
		output, err := p.Text(path)
		if err != nil {
			log.Fatalln(err)
		}

		// Raw
		println("### Raw ###")
		println(output.BodyRaw)

		// Layout
		println("### Layout ###")
		println(output.BodyLayout)
	}
}

var parseImportCmd = &cobra.Command{
	Use:   "import [files]",
	Short: "Import the parsed text from a file into the database",
	Args:  cobra.MinimumNArgs(1),
	Run:   parseImport,
}

func parseImport(cmd *cobra.Command, args []string) {
	a := app.NewApp()

	for _, path := range args {
		err := a.VerificationParse(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
