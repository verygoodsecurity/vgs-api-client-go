package cli

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/verygoodsecurity/vgs-api-client-go/clients"
	"os"
)

func GetOrganizations() {
	organizationClient := clients.NewOrganizationClient(clients.EnvironmentConfig())

	organizations, _ := organizationClient.GetOrganizations()

	renderOrganizationsTable(organizations)
}

func DescribeOrganization(orgId string) {
	organizationClient := clients.NewOrganizationClient(clients.EnvironmentConfig())

	organization, _ := organizationClient.DescribeOrganization(orgId)

	response, _ := json.MarshalIndent(organization, "", "  ")
	fmt.Println(string(response))
}

func renderOrganizationsTable(organizations []clients.Organization) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	for _, org := range organizations {
		table.Append([]string{org.Id, org.Name, org.State, org.CreatedAt})
	}

	table.Render()
}
