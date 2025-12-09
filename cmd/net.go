package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/egichandrap/opsdoctor/internal/netscan"
	"github.com/egichandrap/opsdoctor/utils"
)

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Network diagnostic tools",
}

var netTLSCmd = &cobra.Command{
	Use:   "tls [host]",
	Short: "Check TLS certificate info (host or host:port or https://host)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		utils.VPrintf("tls scan %s\n", host)
		res, err := netscan.CheckTLS(host)
		if err != nil {
			utils.PrintColor("red", err.Error())
			return err
		}
		utils.PrintColor("green", "TLS Scan Result:")
		fmt.Println(res)
		return nil
	},
}

var netPingCmd = &cobra.Command{
	Use:   "ping [host]",
	Short: "Check TCP connectivity (host[:port])",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		utils.VPrintf("ping %s\n", host)
		res := netscan.CheckConnectivity(host)
		fmt.Println(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(netCmd)
	netCmd.AddCommand(netTLSCmd)
	netCmd.AddCommand(netPingCmd)
}
