package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/egichandrap/opsdoctor/internal/apitest"
	"github.com/egichandrap/opsdoctor/utils"
)

var (
	apiMethod      string
	apiBody        string
	apiConcurrency int
	apiRequests    int
	apiTimeout     int
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API diagnostics (single or load test)",
}

var apiTestCmd = &cobra.Command{
	Use:   "test [url]",
	Short: "Test API. Provide one URL or multiple URLs (space separated).",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urls := args
		timeout := time.Duration(apiTimeout) * time.Second

		spec := apitest.NewSpec(strings.ToUpper(apiMethod), urls[0], []byte(apiBody), timeout, nil)

		if len(urls) > 1 && (apiConcurrency > 1 || apiRequests > 1) {
			allSummaries := make(map[string]apitest.Summary)
			for _, u := range urls {
				spec.URL = u
				utils.VPrintf("Running target %s (concurrency=%d, requests=%d)\n", u, apiConcurrency, apiRequests)
				_, summary := apitest.TestConcurrent(spec, apiConcurrency, apiRequests)
				allSummaries[u] = summary
				if utils.IsJSONOutput() {
					b, _ := json.MarshalIndent(summary, "", "  ")
					fmt.Println(string(b))
				} else {
					fmt.Printf("== %s ==\n", u)
					fmt.Printf("Total: %d  Success: %d  Failed: %d  Avg: %v  Min: %v  Max: %v\n",
						summary.TotalRequests, summary.Success, summary.Failed, summary.AvgLatency, summary.MinLatency, summary.MaxLatency)
				}
			}
			if utils.IsJSONOutput() {
				b, _ := json.MarshalIndent(allSummaries, "", "  ")
				fmt.Println(string(b))
			}
			return nil
		}

		if apiConcurrency > 1 || apiRequests > 1 {
			utils.VPrintf("Running load test for %s (concurrency=%d, requests=%d)\n", spec.URL, apiConcurrency, apiRequests)
			_, summary := apitest.TestConcurrent(spec, apiConcurrency, apiRequests)
			if utils.IsJSONOutput() {
				b, _ := apitest.SummaryToJSON(summary)
				fmt.Println(string(b))
				return nil
			}
			fmt.Printf("Total: %d  Success: %d  Failed: %d  Avg: %v  Min: %v  Max: %v  Duration: %v\n",
				summary.TotalRequests, summary.Success, summary.Failed, summary.AvgLatency, summary.MinLatency, summary.MaxLatency, summary.Duration)
			return nil
		}

		res := apitest.TestSingle(spec)
		if utils.IsJSONOutput() {
			b, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(b))
			return nil
		}
		if res.Error != "" {
			utils.PrintColor("red", fmt.Sprintf("%s -> ERROR: %s", res.URL, res.Error))
		} else {
			utils.PrintColor("green", fmt.Sprintf("%s -> %d in %v", res.URL, res.Status, res.Latency))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.AddCommand(apiTestCmd)

	apiTestCmd.Flags().StringVarP(&apiMethod, "method", "m", "GET", "HTTP method")
	apiTestCmd.Flags().StringVarP(&apiBody, "body", "b", "", "HTTP request body (for POST/PUT)")
	apiTestCmd.Flags().IntVarP(&apiConcurrency, "concurrency", "c", 1, "Number of concurrent workers")
	apiTestCmd.Flags().IntVarP(&apiRequests, "requests", "n", 1, "Total requests to send")
	apiTestCmd.Flags().IntVarP(&apiTimeout, "timeout", "t", 10, "Per-request timeout (seconds)")
}
