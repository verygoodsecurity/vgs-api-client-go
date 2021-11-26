package clients

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"os"
	"vgs-api-client/log"
	"vgs-api-client/util"
)

import _ "github.com/joho/godotenv/autoload"

type TenantClient struct {
	accountManagementEndpoint string
	vaultManagementEndpoint   string
	restyClient               resty.Client
	authToken                 string
}

type CreateTenantForm struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

type UpdateTenantForm struct {
	Name string `json:"name"`
}

type Tenant struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	internalId  string
}

func NewTenantClient() *TenantClient {
	restyClient := resty.New()

	return &TenantClient{
		accountManagementEndpoint: os.Getenv("ACCOUNT_MANAGEMENT_API_BASE_URL") + "/vaults",
		vaultManagementEndpoint:   os.Getenv("VAULT_MANAGEMENT_API_BASE_URL") + "/vaults",
		restyClient:               *restyClient,
		authToken:                 util.GetToken(),
	}
}

func (c *TenantClient) request() *resty.Request {
	return c.restyClient.R().
		SetHeader("Accept", "application/vnd.api+json").
		SetHeader("Content-Type", "application/vnd.api+json").
		SetAuthToken(c.authToken)
}

func (c *TenantClient) GetTenants(organizationId string) ([]Tenant, error) {
	tenantsData, _ := c.getTenantsFromAccountManagement()

	var organizationTenants []Tenant
	for _, tenant := range tenantsData.Data {
		if tenant.Relationships.Organization.Data.Id == organizationId {
			organizationTenants = append(organizationTenants, Tenant{
				Id:          tenant.Attributes.Identifier,
				Name:        tenant.Attributes.Name,
				Environment: tenant.Attributes.Environment,
				State:       "-",
				CreatedAt:   tenant.Attributes.CreatedAt,
				UpdatedAt:   tenant.Attributes.UpdatedAt,
			})
		}
	}

	return organizationTenants, nil
}

func (c *TenantClient) Retrieve(tenantId string) (*Tenant, error) {
	tenantsAPIData, _ := c.getTenantsFromAccountManagement()

	var accountManagementTenant tenantAPI

	for _, tnt := range tenantsAPIData.Data {
		if tnt.Attributes.Identifier == tenantId {
			accountManagementTenant = tnt
			break
		}
	}

	tenantData, _ := c.getTenantFromVaultManagement(tenantId)
	vaultManagementTenant := tenantData.Data

	// Unfortunately we have to merge tenant information from two APIs
	var tenant = Tenant{
		Id:          accountManagementTenant.Attributes.Identifier,
		Name:        accountManagementTenant.Attributes.Name,
		Environment: accountManagementTenant.Attributes.Environment,
		State:       vaultManagementTenant.Attributes.State,
		CreatedAt:   accountManagementTenant.Attributes.CreatedAt,
		UpdatedAt:   accountManagementTenant.Attributes.UpdatedAt,
		internalId:  accountManagementTenant.Id,
	}

	return &tenant, nil
}

func (c *TenantClient) SuspendTenant(tenantId string) error {
	tenant, _ := c.Retrieve(tenantId)

	_, err := c.request().Delete(c.accountManagementEndpoint + "/" + tenant.internalId)
	if err != nil {
		return err
	}

	return nil
}

func (c *TenantClient) ProvisionTenant(orgId string, createTenantForm CreateTenantForm) (*Tenant, error) {
	data := tenantAPIData{
		Data: tenantAPI{
			Attributes: tenantAPIAttributes{
				Name:        createTenantForm.Name,
				Environment: createTenantForm.Environment,
			},
			Relationships: tenantAPIRelationships{
				Organization: organizationAPIData{
					Data: organizationAPI{
						Id: orgId,
					},
				},
			},
			TenantType: "vaults",
		},
	}

	payload, err := json.Marshal(data)

	resp, err := c.request().SetBody(payload).Post(c.accountManagementEndpoint)
	if err != nil {
		return nil, err
	}

	var respBody tenantAPIData
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		log.Fatalf("error deserializing data")
	}

	return c.Retrieve(respBody.Data.Attributes.Identifier)
}

func (c *TenantClient) getTenantFromVaultManagement(tenantId string) (*tenantAPIData, error) {
	resp, err := c.request().SetHeader("VGS-Tenant", tenantId).Get(c.vaultManagementEndpoint + "/" + tenantId)
	if err != nil {
		return nil, err
	}

	var tenantAPIData tenantAPIData
	if err := json.Unmarshal(resp.Body(), &tenantAPIData); err != nil {
		log.Fatalf("error deserializing data")
	}

	return &tenantAPIData, nil
}

func (c *TenantClient) getTenantsFromAccountManagement() (*tenantsAPIData, error) {
	resp, err := c.request().Get(c.accountManagementEndpoint)
	if err != nil {
		return nil, err
	}

	var tenantsAPIData tenantsAPIData
	if err := json.Unmarshal(resp.Body(), &tenantsAPIData); err != nil {
		log.Fatalf("error deserializing data")
	}

	return &tenantsAPIData, nil
}

// For internal use

type tenantsAPIData struct {
	Data []tenantAPI `json:"data,omitempty"`
}

type tenantAPIData struct {
	Data tenantAPI `json:"data,omitempty"`
}

type tenantAPI struct {
	Id            string                 `json:"id,omitempty"`
	Identifier    string                 `json:"identifier,omitempty"`
	TenantType    string                 `json:"type,omitempty"`
	Links         tenantAPILinks         `json:"links,omitempty"`
	Relationships tenantAPIRelationships `json:"relationships,omitempty"`
	Credentials   tenantAPICredentials   `json:"credentials,omitempty"`
	Attributes    tenantAPIAttributes    `json:"attributes,omitempty"`
}

type tenantAPILinks struct {
	ForwardProxy       string `json:"forward_proxy,omitempty"`
	ReverseProxy       string `json:"reverse_proxy,omitempty"`
	VaultApi           string `json:"vault_api,omitempty"`
	VaultManagementApi string `json:"vault_management_api,omitempty"`
}

type tenantAPICredentials struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type tenantAPIRelationships struct {
	Organization organizationAPIData `json:"organization,omitempty"`
}

type tenantAPIAttributes struct {
	Id          string   `json:"id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Name        string   `json:"name,omitempty"`
	Environment string   `json:"environment,omitempty"`
	State       string   `json:"state,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}
