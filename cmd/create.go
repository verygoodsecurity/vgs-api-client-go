package cmd

import (
	"github.com/spf13/cobra"
	"vgs-api-client/cli"
)

import _ "github.com/joho/godotenv/autoload"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create [tenant]",
}

func init() {
	var tenantName string
	var tenantEnvironment string

	var createTenantCmd = &cobra.Command{
		Use:   "tenant [account id]",
		Short: "Create tenant",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.ProvisionTenant(args[0], tenantName, tenantEnvironment)
		},
	}

	rootCmd.AddCommand(createCmd)
	createTenantCmd.Flags().StringVarP(&tenantName, "tenant", "t", "", "Tenant name (required)")
	createTenantCmd.Flags().StringVarP(&tenantEnvironment, "environment", "e", "SANDBOX", "Tenant environment")
	createTenantCmd.MarkFlagRequired("tenant")
	createCmd.AddCommand(createTenantCmd)
}
