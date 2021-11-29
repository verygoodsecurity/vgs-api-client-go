package clients

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/verygoodsecurity/vgs-api-client-go/tools"
	"io"
	"strings"
)

type RouteClient struct {
	apiBase     string
	restyClient resty.Client
	authToken   string
}

type idHolder struct {
	Id string `json:"id"`
}

func (r *RouteClient) ImportRoute(tenant string, vgsYaml io.Reader) error {
	yaml, err := reader2string(vgsYaml)
	if err != nil {
		return errors.Wrap(err, "failed to read YAML")
	}
	routeJson, err := tools.Yaml2Json(yaml)
	if err != nil {
		return errors.Wrap(err, "failed to convert YAML to JSON")
	}

	var holder idHolder
	err = json.Unmarshal([]byte(routeJson), &holder)
	if err != nil || holder.Id == "" {
		return errors.Wrap(err, "failed to extract ID from route")
	}

	response, err := r.request().
		SetHeader("VGS-Tenant", tenant).
		SetBody(routeJson).
		Put(fmt.Sprintf("%s/rule-chains/%s", r.apiBase, holder.Id))
	if err != nil {
		return errors.Wrap(err, "API request failed")
	}
	_ = response
	return nil
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
