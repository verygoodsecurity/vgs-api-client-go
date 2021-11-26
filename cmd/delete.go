package cmd

import (
	"github.com/spf13/cobra"
	"github.com/verygoodsecurity/vgs-api-client-go/cli"
)

import _ "github.com/joho/godotenv/autoload"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete [tenant]",
}

var deleteTenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "Delete tenant",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.SuspendTenant(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteTenantCmd)
}
