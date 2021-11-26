package cli

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/verygoodsecurity/vgs-api-client-go/clients"
	"os"
)

func GetVaults(organizationId string) {
	tenantClient := clients.NewVaultClient()

	tenants, _ := tenantClient.GetVaults(organizationId)

	renderVaultsTable(tenants)
}

func GetVault(vaultId string) {
	tenantClient := clients.NewVaultClient()

	tenant, _ := tenantClient.RetrieveVault(vaultId)

	response, _ := json.MarshalIndent(tenant, "", "  ")
	fmt.Println(string(response))
}

func ProvisionVault(orgId string, name string, environment string) {
	tenantClient := clients.NewVaultClient()

	tenant, _ := tenantClient.ProvisionVault(
		orgId,
		clients.CreateVaultForm{
			Name:        name,
			Environment: environment,
		})

	GetVault(tenant.Id)
}

func SuspendVault(tenantId string) {
	tenantClient := clients.NewVaultClient()

	err := tenantClient.SuspendVault(tenantId)
	if err != nil {
		fmt.Println("Failed to delete tenant", err)
		return
	}

	fmt.Println("Deleted tenant " + tenantId)
}

func renderVaultsTable(tenants []clients.Vault) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for _, tenant := range tenants {
		table.Append([]string{
			tenant.Id,
			tenant.Name,
			tenant.Environment,
			tenant.CreatedAt})
	}

	table.Render()
}
