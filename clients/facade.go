package clients

type VgsClientsFacade interface {
	Organizations(ClientConfig) *OrganizationClient
	Vaults(ClientConfig) *VaultClient
	Routes(ClientConfig) *RouteClient
}

type facade struct{}

func (f *facade) Organizations(config ClientConfig) *OrganizationClient {
	return NewOrganizationClient(config)
}

func (f *facade) Vaults(config ClientConfig) *VaultClient {
	return NewVaultClient(config)
}

func (f *facade) Routes(config ClientConfig) *RouteClient {
	return NewRouteClient(config)
}

func NewVgsFacade() VgsClientsFacade {
	return &facade{}
}
