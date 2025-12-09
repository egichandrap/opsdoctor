package cmd

import (
	"fmt"
	"github.com/egichandrap/opsdoctor/internal/springcheck"
	"github.com/egichandrap/opsdoctor/utils"
	"github.com/spf13/cobra"
)

var springCmd = &cobra.Command{
	Use:   "spring",
	Short: "Spring Boot diagnostics",
}

var springAnalyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze Spring Boot log",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		fmt.Println(utils.Green("=== Spring Log Analyzer ==="))

		res, err := springcheck.Analyze(file)
		if err != nil {
			fmt.Println(utils.Red("Error:"), err)
			return
		}

		fmt.Println(utils.Yellow("\n--- Errors ---"))
		for _, e := range res.Errors {
			fmt.Println(e)
		}

		fmt.Println(utils.Yellow("\n--- Slow Queries ---"))
		for _, q := range res.SlowQueries {
			fmt.Println(q)
		}

		fmt.Printf("\nFound %d errors, %d slow queries\n", len(res.Errors), len(res.SlowQueries))
	},
}

func init() {
	springCmd.AddCommand(springAnalyzeCmd)
	rootCmd.AddCommand(springCmd)
}
