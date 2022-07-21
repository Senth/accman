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

func initSie(rootCmd *cobra.Command) {
	rootCmd.AddCommand(sieCmd)
	sieCmd.AddCommand(sieImportCmd)
	//sieCmd.AddCommand(sieExportCmd)
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
