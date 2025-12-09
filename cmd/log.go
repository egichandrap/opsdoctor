package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/egichandrap/opsdoctor/internal/loganalyzer"
	"github.com/egichandrap/opsdoctor/utils"
)

var (
	logThreshold int
	logPattern   string
	logExport    string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log analysis utilities",
}

var logAnalyzeCmd = &cobra.Command{
	Use:   "analyze [file]",
	Short: "Analyze log file for slow entries (e.g. '123ms')",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		utils.VPrintf("Starting log analyze: %s (threshold=%d)\n", path, logThreshold)

		result, txt, err := loganalyzer.Analyze(path, logThreshold, logPattern, logExport)
		if err != nil {
			return err
		}

		if utils.IsJSONOutput() {
			b, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		}

		fmt.Println(txt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.AddCommand(logAnalyzeCmd)

	logAnalyzeCmd.Flags().IntVarP(&logThreshold, "slow", "s", 1000, "threshold in ms for slow logs")
	logAnalyzeCmd.Flags().StringVarP(&logPattern, "pattern", "p", "", "optional substring filter for log lines")
	logAnalyzeCmd.Flags().StringVarP(&logExport, "export", "e", "", "optional JSON export path (writes analysis JSON)")
}
