package clients

import (
	"github.com/Nerzal/gocloak/v8"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// private interface for internal usage
type authServerClient interface {
	GetToken() (string, error)
}

type keycloak struct {
	config ClientConfig
}

func newKeycloak(config ClientConfig) *keycloak {
	return &keycloak{
		config: config,
	}
}

func (a *keycloak) GetToken() (string, error) {
	client := gocloak.NewClient(a.config.Get("VGS_KEYCLOAK_URL"))
	ctx := context.Background()
	token, err := client.GetToken(ctx,
		a.config.Get("VGS_KEYCLOAK_REALM"),
		gocloak.TokenOptions{
			ClientID:     strptr(a.config.Get("VGS_CLIENT_ID")),
			ClientSecret: strptr(a.config.Get("VGS_CLIENT_SECRET")),
			GrantType:    strptr("client_credentials"),
		})
	if err != nil {
		return "", errors.Wrap(err, "failed to authenticate")
	}
	return token.AccessToken, nil
}

func strptr(str string) *string {
	return &str
}
