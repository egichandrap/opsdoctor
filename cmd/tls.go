package cmd

import (
	"fmt"

	"github.com/egichandrap/opsdoctor/internal/tlscheck"
	"github.com/egichandrap/opsdoctor/utils"
	"github.com/spf13/cobra"
)

var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "Standalone TLS commands",
}

var tlsCheckCmd = &cobra.Command{
	Use:   "check [host]",
	Short: "Check TLS certificate info (alias to net tls)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		utils.VPrintf("tls check %s\n", host)
		res, err := tlscheck.CheckTLS(host)
		if err != nil {
			utils.PrintColor("red", err.Error())
			return err
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
	tlsCmd.AddCommand(tlsCheckCmd)
}
