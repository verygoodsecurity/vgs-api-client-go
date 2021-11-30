package clients

import "github.com/go-resty/resty/v2"

type client struct {
	rest *resty.Client
	auth authServerClient

	cachedToken string
}

func (c client) request() (*resty.Request, error) {
	if c.cachedToken == "" {
		token, err := c.auth.GetToken()
		if err != nil {
			return nil, err
		}
		c.cachedToken = token
	}
	return c.rest.R().
		SetHeader("Accept", "application/vnd.api+json").
		SetHeader("Content-Type", "application/vnd.api+json").
		SetAuthToken(c.cachedToken), nil
}
