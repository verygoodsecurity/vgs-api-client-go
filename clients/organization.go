package clients

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

import _ "github.com/joho/godotenv/autoload"

type Organization struct {
	Id           string                    `json:"id"`
	Name         string                    `json:"name"`
	State        string                    `json:"state"`
	CreatedAt    string                    `json:"created_at"`
	UpdatedAt    string                    `json:"updated_at"`
	Users        []OrganizationUser        `json:"users"`
	Environments []OrganizationEnvironment `json:"environments"`
}

type OrganizationUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrganizationEnvironment struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Region     string `json:"region"`
}

type OrganizationClient struct {
	endpoint string
	client
}

func NewOrganizationClient(config ClientConfig) *OrganizationClient {
	return &OrganizationClient{
		endpoint: config.Get("ACCOUNT_MANAGEMENT_API_BASE_URL") + "/organizations",
		client: client{
			rest: resty.New(),
			auth: newKeycloak(config),
		},
	}
}

func (c *OrganizationClient) GetOrganizations() ([]Organization, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.Get(c.endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var organizationsAPIData organizationsAPIData
	if err := json.Unmarshal(resp.Body(), &organizationsAPIData); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	var organizations []Organization
	for _, org := range organizationsAPIData.Data {
		organizations = append(organizations, Organization{
			Id:        org.Id,
			Name:      org.Attributes.Name,
			State:     org.Attributes.State,
			CreatedAt: org.Attributes.CreatedAt,
			UpdatedAt: org.Attributes.UpdatedAt,
		})
	}

	return organizations, nil
}

func (c *OrganizationClient) DescribeOrganization(orgId string) (*Organization, error) {
	users, _ := c.getOrganizationUsers(orgId)
	environments, _ := c.getOrganizationEnvironments(orgId)
	organization, _ := c.getOrganization(orgId)

	organization.Users = users
	organization.Environments = environments

	return organization, nil
}

func (c *OrganizationClient) getOrganization(orgId string) (*Organization, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.Get(c.endpoint + "/" + orgId)
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var organizationAPIData organizationAPIData
	if err := json.Unmarshal(resp.Body(), &organizationAPIData); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	organization := Organization{
		Id:        organizationAPIData.Data.Id,
		Name:      organizationAPIData.Data.Attributes.Name,
		State:     organizationAPIData.Data.Attributes.State,
		CreatedAt: organizationAPIData.Data.Attributes.CreatedAt,
		UpdatedAt: organizationAPIData.Data.Attributes.UpdatedAt,
	}

	return &organization, nil
}

func (c *OrganizationClient) getOrganizationEnvironments(orgId string) ([]OrganizationEnvironment, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.Get(c.endpoint + "/" + orgId + "/environments")
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var organizationEnvironmentsAPIData organizationEnvironmentsAPIData
	if err := json.Unmarshal(resp.Body(), &organizationEnvironmentsAPIData); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	var environments []OrganizationEnvironment
	for _, environmentAPI := range organizationEnvironmentsAPIData.Data {
		environments = append(environments, OrganizationEnvironment{
			Id:         environmentAPI.Id,
			Name:       environmentAPI.Attributes.Name,
			Identifier: environmentAPI.Attributes.Identifier,
			Region:     environmentAPI.Attributes.Region,
		})
	}

	return environments, nil
}

func (c *OrganizationClient) getOrganizationUsers(orgId string) ([]OrganizationUser, error) {
	request, err := c.client.request()
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}
	resp, err := request.Get(c.endpoint + "/" + orgId + "/users")
	if err != nil {
		return nil, errors.Wrap(err, "API request failed")
	}

	var organizationUsersAPIData organizationUsersAPIData
	if err := json.Unmarshal(resp.Body(), &organizationUsersAPIData); err != nil {
		return nil, errors.Wrap(err, "error deserializing data")
	}

	var users []OrganizationUser
	for _, userAPI := range organizationUsersAPIData.Data {
		users = append(users, OrganizationUser{
			Id:        userAPI.Id,
			Name:      userAPI.Attributes.Name,
			Email:     userAPI.Attributes.EmailAddress,
			CreatedAt: userAPI.Attributes.CreatedAt,
			UpdatedAt: userAPI.Attributes.UpdatedAt,
		})
	}

	return users, nil
}

// For internal use

type organizationsAPIData struct {
	Data []organizationAPI `json:"data,omitempty"`
}

type organizationAPIData struct {
	Data organizationAPI `json:"data,omitempty"`
}

type organizationAPI struct {
	Id         string                    `json:"id,omitempty"`
	OrgType    string                    `json:"type,omitempty"`
	Attributes organizationAPIAttributes `json:"attributes,omitempty"`
}

type organizationAPIAttributes struct {
	InternalId  string   `json:"internal_id,omitempty"`
	Identifier  string   `json:"identifier,omitempty"`
	Name        string   `json:"name,omitempty"`
	Active      bool     `json:"active,omitempty"`
	State       string   `json:"state,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type organizationEnvironmentsAPIData struct {
	Data []organizationEnvironmentAPI `json:"data,omitempty"`
}

type organizationEnvironmentAPI struct {
	Id         string                               `json:"id,omitempty"`
	Attributes organizationEnvironmentAPIAttributes `json:"attributes,omitempty"`
}

type organizationEnvironmentAPIAttributes struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Region     string `json:"region,omitempty"`
	IsLive     bool   `json:"is_live,omitempty"`
}

type organizationUsersAPIData struct {
	Data []organizationUserAPI `json:"data,omitempty"`
}

type organizationUserAPI struct {
	Id         string                        `json:"id,omitempty"`
	Attributes organizationUserAPIAttributes `json:"attributes,omitempty"`
}

type organizationUserAPIAttributes struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}
