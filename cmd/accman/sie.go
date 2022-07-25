package main

import (
	"github.com/Senth/accman/internal/app"
	"github.com/spf13/cobra"
	"log"
)

var sieCmd = &cobra.Command{
	Use:   "sie",
	Short: "Import and export SIE files",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var outputPath string

func initSie(rootCmd *cobra.Command) {
	rootCmd.AddCommand(sieCmd)
	sieCmd.AddCommand(sieImportCmd)
	sieCmd.AddCommand(sieExportCmd)

	sieExportCmd.Flags().StringVar(&outputPath, "output", "", "Output path instead of default 'year.se'")
}

var sieImportCmd = &cobra.Command{
	Use:   "import [files]",
	Short: "Import SIE files",
	Args:  cobra.MinimumNArgs(1),
	Run:   sieImport,
}

func sieImport(cmd *cobra.Command, args []string) {
	a := app.NewApp()

	for _, path := range args {
		err := a.SIEImport(path)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

var sieExportCmd = &cobra.Command{
	Use:   "export year",
	Short: "Export a fiscal year to a SIE file",
	Args:  cobra.ExactArgs(1),
	Run:   sieExport,
}

func sieExport(cmd *cobra.Command, args []string) {
	a := app.NewApp()

	err := a.SIEExport(args[0], outputPath)
	if err != nil {
		log.Fatalln(err)
	}
}
