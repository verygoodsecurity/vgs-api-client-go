package cmd

import (
	"github.com/spf13/cobra"
	"vgs-api-client/cli"
)

import _ "github.com/joho/godotenv/autoload"

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe [organization|tenant]",
}

var describeTenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Describe tenant",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetTenant(args[0])
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
	describeCmd.AddCommand(describeTenantCmd)
	describeCmd.AddCommand(describeOrganizationCmd)
}
