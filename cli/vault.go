package cli

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/verygoodsecurity/vgs-api-client-go/clients"
	"os"
)

func GetVaults(organizationId string) {
	tenantClient := clients.NewVaultClient(clients.EnvironmentConfig())

	tenants, _ := tenantClient.GetVaults(organizationId)

	renderVaultsTable(tenants)
}

func GetVault(vaultId string) {
	vaultClient := clients.NewVaultClient(clients.EnvironmentConfig())

	tenant, _ := vaultClient.RetrieveVault(vaultId)

	response, _ := json.MarshalIndent(tenant, "", "  ")
	fmt.Println(string(response))
}

func ProvisionVault(orgId string, name string, environment string) {
	vaultClient := clients.NewVaultClient(clients.EnvironmentConfig())

	tenant, _ := vaultClient.ProvisionVault(
		orgId,
		clients.CreateVaultForm{
			Name:        name,
			Environment: environment,
		})

	GetVault(tenant.Id)
}

func SuspendVault(tenantId string) {
	vaultClient := clients.NewVaultClient(clients.EnvironmentConfig())

	err := vaultClient.SuspendVault(tenantId)
	if err != nil {
		fmt.Println("Failed to delete tenant", err)
		return
	}

	fmt.Println("Deleted tenant " + tenantId)
}

func renderVaultsTable(vaults []clients.Vault) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for _, vault := range vaults {
		table.Append([]string{
			vault.Id,
			vault.Name,
			vault.Environment,
			vault.CreatedAt})
	}

	table.Render()
}
