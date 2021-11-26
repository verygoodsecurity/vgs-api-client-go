package cmd

import (
	"github.com/spf13/cobra"
	"vgs-api-client/cli"
)

import _ "github.com/joho/godotenv/autoload"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get [organizations|tenants]",
}

var getOrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "Get organizations",
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetOrganizations()
	},
}

var getTenantsCmd = &cobra.Command{
	Use:   "tenants [org_id]",
	Short: "Get tenants for organization",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetTenants(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getOrganizationsCmd)
	getCmd.AddCommand(getTenantsCmd)
}
