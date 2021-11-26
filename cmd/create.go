package cmd

import (
	"github.com/spf13/cobra"
	"github.com/verygoodsecurity/vgs-api-client-go/cli"
)

import _ "github.com/joho/godotenv/autoload"

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create [vault]",
}

func init() {
	var vaultName string
	var vaultEnvironment string

	var createTenantCmd = &cobra.Command{
		Use:   "vault [account id]",
		Short: "Create vault",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.ProvisionVault(args[0], vaultName, vaultEnvironment)
		},
	}

	rootCmd.AddCommand(createCmd)
	createTenantCmd.Flags().StringVarP(&vaultName, "vault", "v", "", "Vault name (required)")
	createTenantCmd.Flags().StringVarP(&vaultEnvironment, "environment", "e", "SANDBOX", "Vault environment")
	createTenantCmd.MarkFlagRequired("vault")
	createCmd.AddCommand(createTenantCmd)
}
