package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/egichandrap/opsdoctor/utils"
)

var (
	NoColor   bool
	Verbose   bool
	OutFormat string // "text" or "json"
)

var rootCmd = &cobra.Command{
	Use:   "opsdoctor",
	Short: "OpsDoctor - RHEL & Spring Diagnostic Toolkit",
	Long:  "OpsDoctor is a CLI toolkit for diagnostics (network, TLS, logs, API, services).",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&OutFormat, "output", "o", "text", "output format: text|json")

	cobra.OnInitialize(func() {
		utils.SetNoColor(NoColor)
		utils.SetVerbose(Verbose)
		utils.SetOutFormat(OutFormat)
	})
}
