package clients

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/verygoodsecurity/vgs-api-client-go/tools"
	"io"
	"strings"
)

type RouteClient struct {
	apiBase     string
	restyClient *resty.Client
	authToken   string
}

func NewRouteClient(config ClientConfig) *RouteClient {
	return &RouteClient{
		apiBase:     config.Get("VGS_VAULT_MANAGEMENT_API_BASE_URL"),
		restyClient: resty.New(),
		authToken:   newKeycloak(config).GetToken(),
	}
}

func (r *RouteClient) ImportRoute(vault string, vgsYaml io.Reader) error {
	yaml, err := reader2string(vgsYaml)
	if err != nil {
		return errors.Wrap(err, "failed to read YAML")
	}
	id, err := tools.RouteIdFromYaml(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to extract ID from route")
	}

	routeJson, err := tools.Yaml2Json(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to convert YAML to JSON")
	}

	_, err = r.request().
		SetHeader("VGS-Tenant", vault).
		SetBody(routeJson).
		Put(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, id))
	return errors.Wrap(err, "API request failed")
}

func (r *RouteClient) DeleteRoute(vault, id string) error {
	_, err := r.request().
		SetHeader("VGS-Tenant", vault).
		Delete(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, id))
	return errors.Wrap(err, "API request failed")
}

func (r *RouteClient) request() *resty.Request {
	return r.restyClient.R().
		SetHeader("Accept", "application/vnd.api+json").
		SetHeader("Content-Type", "application/vnd.api+json").
		SetAuthToken(r.authToken)
}

func reader2string(reader io.Reader) (string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
