package clients

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

import _ "github.com/joho/godotenv/autoload"

type VaultClient struct {
	accountManagementEndpoint string
	vaultManagementEndpoint   string
	client
}

type CreateVaultForm struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
}

type UpdateVaultForm struct {
	Name string `json:"name"`
}

type Vault struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	State       string `json:"state"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	internalId  string
}

func NewVaultClient(config ClientConfig) *VaultClient {
	return &VaultClient{
		accountManagementEndpoint: config.Get("VGS_ACCOUNT_MANAGEMENT_API_BASE_URL") + "/vaults",
		vaultManagementEndpoint:   config.Get("VGS_VAULT_MANAGEMENT_API_BASE_URL") + "/vaults",
		client: client{
			rest: resty.New(),
			auth: newKeycloak(config),
		},
	}
}

func (c *VaultClient) GetVaults(organizationId string) ([]Vault, error) {
	vaultsData, _ := c.getVaultsFromAccountManagement()

	var organizationVaults []Vault
	for _, vault := range vaultsData.Data {
		if vault.Relationships.Organization.Data.Id == organizationId {
			organizationVaults = append(organizationVaults, Vault{
				Id:          vault.Attributes.Identifier,
				Name:        vault.Attributes.Name,
				Environment: vault.Attributes.Environment,
				State:       "-",
				CreatedAt:   vault.Attributes.CreatedAt,
				UpdatedAt:   vault.Attributes.UpdatedAt,
			})
		}
	}

	return organizationVaults, nil
}

func (c *VaultClient) RetrieveVault(vaultId string) (*Vault, error) {
	vaultsAPIData, _ := c.getVaultsFromAccountManagement()

	var accountManagementVault vaultAPI

	for _, tnt := range vaultsAPIData.Data {
		if tnt.Attributes.Identifier == vaultId {
			accountManagementVault = tnt
			break
		}
	}

	vaultData, _ := c.getVaultFromVaultManagement(vaultId)
	vaultManagementVault := vaultData.Data

	// Unfortunately we have to merge vault information from two APIs
	var vault = Vault{
		Id:          accountManagementVault.Attributes.Identifier,
		Name:        accountManagementVault.Attributes.Name,
		Environment: accountManagementVault.Attributes.Environment,
		State:       vaultManagementVault.Attributes.State,
		CreatedAt:   accountManagementVault.Attributes.CreatedAt,
		UpdatedAt:   accountManagementVault.Attributes.UpdatedAt,
		internalId:  accountManagementVault.Id,
	}

	return &vault, nil
}

func (c *VaultClient) SuspendVault(vaultId string) error {
	vault, err := c.RetrieveVault(vaultId)
	if err != nil {
		return errors.Wrap(err, "failed to find vault by ID")
	}
	request, err := c.client.request()
	if err != nil {
		return errors.Wrap(err, "API request failed")
	}
	_, err = request.Delete(c.accountManagementEndpoint + "/" + vault.internalId)
	return errors.Wrap(err, "API request failed")
}

func (c *VaultClient) ProvisionVault(orgId string, createVaultForm CreateVaultForm) (*Vault, error) {
	data := vaultAPIData{
		Data: vaultAPI{
			Attributes: vaultAPIAttributes{
				Name:        createVaultForm.Name,
				Environment: createVaultForm.Environment,
			},
			Relationships: vaultAPIRelationships{
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
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.SetBody(payload).Post(c.accountManagementEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var respBody vaultAPIData
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	return c.RetrieveVault(respBody.Data.Attributes.Identifier)
}

func (c *VaultClient) getVaultFromVaultManagement(vaultId string) (*vaultAPIData, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.SetHeader("VGS-Tenant", vaultId).Get(c.vaultManagementEndpoint + "/" + vaultId)
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var data vaultAPIData
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	return &data, nil
}

func (c *VaultClient) getVaultsFromAccountManagement() (*vaultsAPIData, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.Get(c.accountManagementEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var data vaultsAPIData
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}
	return &data, nil
}

// For internal use

type vaultsAPIData struct {
	Data []vaultAPI `json:"data,omitempty"`
}

type vaultAPIData struct {
	Data vaultAPI `json:"data,omitempty"`
}

type vaultAPI struct {
	Id            string                `json:"id,omitempty"`
	Identifier    string                `json:"identifier,omitempty"`
	TenantType    string                `json:"type,omitempty"`
	Links         vaultAPILinks         `json:"links,omitempty"`
	Relationships vaultAPIRelationships `json:"relationships,omitempty"`
	Credentials   vaultAPICredentials   `json:"credentials,omitempty"`
	Attributes    vaultAPIAttributes    `json:"attributes,omitempty"`
}

type vaultAPILinks struct {
	ForwardProxy       string `json:"forward_proxy,omitempty"`
	ReverseProxy       string `json:"reverse_proxy,omitempty"`
	VaultApi           string `json:"vault_api,omitempty"`
	VaultManagementApi string `json:"vault_management_api,omitempty"`
}

type vaultAPICredentials struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type vaultAPIRelationships struct {
	Organization organizationAPIData `json:"organization,omitempty"`
}

type vaultAPIAttributes struct {
	Id          string   `json:"id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Name        string   `json:"name,omitempty"`
	Environment string   `json:"environment,omitempty"`
	State       string   `json:"state,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}
