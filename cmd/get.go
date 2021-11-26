package cmd

import (
	"github.com/spf13/cobra"
	"github.com/verygoodsecurity/vgs-api-client-go/cli"
)

import _ "github.com/joho/godotenv/autoload"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get [organizations|vaults]",
}

var getOrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "Get organizations",
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetOrganizations()
	},
}

var getTenantsCmd = &cobra.Command{
	Use:   "vaults [org_id]",
	Short: "Get vaults for organization",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.GetVaults(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getOrganizationsCmd)
	getCmd.AddCommand(getTenantsCmd)
}
