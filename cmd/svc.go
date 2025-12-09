package cmd

import (
	"fmt"
	"github.com/egichandrap/opsdoctor/internal/svcchecker"
	"github.com/egichandrap/opsdoctor/utils"
	"github.com/spf13/cobra"
)

var svcCmd = &cobra.Command{
	Use:   "svc",
	Short: "Service diagnostics",
}

var svcStatusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Check systemd service status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		svc := args[0]
		fmt.Println(utils.Green("=== Service Status ==="))
		fmt.Println("Service:", svc)

		out, err := svcchecker.CheckService(svc)
		if err != nil {
			fmt.Println(utils.Red("ERROR:"), err)
		}

		fmt.Println("\n" + out)
	},
}

func init() {
	svcCmd.AddCommand(svcStatusCmd)
	rootCmd.AddCommand(svcCmd)
}
