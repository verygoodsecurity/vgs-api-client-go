package cli

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"vgs-api-client/clients"
)

func GetTenants(organizationId string) {
	tenantClient := clients.NewTenantClient()

	tenants, _ := tenantClient.GetTenants(organizationId)

	renderTenantsTable(tenants)
}

func GetTenant(tenantId string) {
	tenantClient := clients.NewTenantClient()

	tenant, _ := tenantClient.Retrieve(tenantId)

	response, _ := json.MarshalIndent(tenant, "", "  ")
	fmt.Println(string(response))
}

func ProvisionTenant(orgId string, name string, environment string) {
	tenantClient := clients.NewTenantClient()

	tenant, _ := tenantClient.ProvisionTenant(
		orgId,
		clients.CreateTenantForm{
			Name:        name,
			Environment: environment,
		})

	GetTenant(tenant.Id)
}

func SuspendTenant(tenantId string) {
	tenantClient := clients.NewTenantClient()

	err := tenantClient.SuspendTenant(tenantId)
	if err != nil {
		fmt.Println("Failed to delete tenant", err)
		return
	}

	fmt.Println("Deleted tenant " + tenantId)
}

func renderTenantsTable(tenants []clients.Tenant) {
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
