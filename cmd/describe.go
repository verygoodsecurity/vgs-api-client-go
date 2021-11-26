package cmd

import (
	"github.com/spf13/cobra"
	"github.com/verygoodsecurity/vgs-api-client-go/cli"
)

import _ "github.com/joho/godotenv/autoload"

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe [organization|vault]",
}

var describeVaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Describe vault",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetVault(args[0])
	},
}

var describeOrganizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Describe organization",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.DescribeOrganization(args[0])
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.AddCommand(describeVaultCmd)
	describeCmd.AddCommand(describeOrganizationCmd)
}
