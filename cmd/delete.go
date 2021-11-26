package cmd

import (
	"github.com/spf13/cobra"
	"github.com/verygoodsecurity/vgs-api-client-go/cli"
)

import _ "github.com/joho/godotenv/autoload"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete [vault]",
}

var deleteTenantCmd = &cobra.Command{
	Use:   "vault",
	Short: "Delete vault",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.SuspendVault(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteTenantCmd)
}
