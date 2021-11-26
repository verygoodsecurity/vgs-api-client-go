package clients

type VgsClientsFacade interface {
	Organizations() *OrganizationClient
	Tenants() *VaultClient
}

type facade struct {
	organizationClient *OrganizationClient
	vaultsClient       *VaultClient
}

func (f *facade) Organizations() *OrganizationClient {
	return f.organizationClient
}

func (f *facade) Tenants() *VaultClient {
	return f.vaultsClient
}

func NewVgsFacade(config ClientConfig) VgsClientsFacade {
	return &facade{
		organizationClient: NewOrganizationClient(config),
		vaultsClient:       NewVaultClient(config),
	}
}
