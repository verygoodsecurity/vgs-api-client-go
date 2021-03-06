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
	apiBase string
	client
}

func NewRouteClient(config ClientConfig) *RouteClient {
	return &RouteClient{
		apiBase: config.Get("VGS_VAULT_MANAGEMENT_API_BASE_URL"),
		client: client{
			rest: resty.New(),
			auth: newKeycloak(config),
		},
	}
}

func (r *RouteClient) GetRoute(vault, routeId string) (routeJson string, err error) {
	// TODO need to return a struct but it's ok for now
	request, err := r.client.request()
	if err != nil {
		return "", errors.Wrap(err, "API request failed")
	}
	response, err := request.
		SetHeader("VGS-Tenant", vault).
		Get(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, routeId))
	if err != nil {
		return "", errors.Wrap(err, "API request failed")
	}
	if response.StatusCode() == 404 {
		return "", errors.New("Route not found")
	}
	return string(response.Body()), errors.Wrap(err, "API request failed")
}

func (r *RouteClient) ImportRoute(vault string, vgsYaml io.Reader) (id string, err error) {
	yaml, err := reader2string(vgsYaml)
	if err != nil {
		return "", errors.Wrap(err, "failed to read YAML")
	}
	id, err = tools.RouteIdFromYaml(yaml)
	if err != nil {
		return "", errors.Wrap(err, "failed to extract ID from route")
	}

	routeJson, err := tools.Yaml2Json(yaml)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert YAML to JSON")
	}

	requestBody, err := tools.WrapJSON("data", routeJson)
	if err != nil {
		return "", errors.Wrap(err, "failed to construct request body")
	}
	request, err := r.client.request()
	if err != nil {
		return "", errors.Wrap(err, "API request failed")
	}
	response, err := request.
		SetHeader("VGS-Tenant", vault).
		SetBody(requestBody).
		Put(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, id))
	if response.StatusCode() != 200 && response.StatusCode() != 201 {
		return "", errors.Errorf("API returned %s", response.Status())
	}

	return id, errors.Wrap(err, "API request failed")
}

func (r *RouteClient) DeleteRoute(vault, id string) error {
	request, err := r.client.request()
	if err != nil {
		return errors.Wrap(err, "API request failed")
	}
	_, err = request.
		SetHeader("VGS-Tenant", vault).
		Delete(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, id))
	return errors.Wrap(err, "API request failed")
}

func reader2string(reader io.Reader) (string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
