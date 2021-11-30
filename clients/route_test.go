package clients

import (
	asserting "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Set dev VGS_CLIENT_ID/VGS_CLIENT_SECRET in env to run

var routeClient = NewRouteClient(DynamicConfig().
	WithFallback(EnvironmentConfig()).
	AddParameter("VGS_VAULT_MANAGEMENT_API_BASE_URL", "https://api.verygoodsecurity.io").
	AddParameter("VGS_KEYCLOAK_URL", "https://auth.verygoodsecurity.io").
	AddParameter("VGS_KEYCLOAK_REALM", "vgs"))

func TestImportRoute(t *testing.T) {
	assert := asserting.New(t)
	routeId, err := routeClient.ImportRoute("tntbcduzut5", strings.
		NewReader(`id: 04b2e1b7-fb60-472f-a79f-af7e2353f122
type: rule_chain
attributes:
  created_at: '2021-11-26T18:10:08'
  destination_override_endpoint: 'https://echo.apps.verygood.systems'
  entries:
    - classifiers: {}
      config:
        condition: AND
        rules:
          - expression:
              field: PathInfo
              operator: matches
              type: string
              values:
                - /post
          - expression:
              field: ContentType
              operator: equals
              type: string
              values:
                - application/json
            rules: null
      id: 39f2f5db-06a0-461d-9387-dd9a7ab19035
      id_selector: null
      operation: REDACT
      operations: null
      phase: REQUEST
      public_token_generator: UUID
      targets:
        - body
      token_manager: PERSISTENT
      transformer: JSON_PATH
      transformer_config:
        - $.account_number
      transformer_config_map: null
  host_endpoint: (.*)\.verygoodproxy\.io
  id: 04b2e1b7-fb60-472f-a79f-af7e2353f122
  ordinal: null
  port: 80
  protocol: http
  source_endpoint: '*'
  tags:
    name: echo.apps.verygood.systems-beige-crescent
    source: RouteContainer
  updated_at: '2021-11-26T18:10:08'`))
	assert.Nil(err)
	assert.Equal("04b2e1b7-fb60-472f-a79f-af7e2353f122", routeId)
	route, err := routeClient.GetRoute("tntbcduzut5", "04b2e1b7-fb60-472f-a79f-af7e2353f122")
	assert.Nil(err)
	assert.NotEqual("", route)
}

func TestRouteGet_NotExisting(t *testing.T) {
	assert := asserting.New(t)
	route, err := routeClient.GetRoute("tntbcduzut5", "f954c07f-4d5f-4a3a-b620-d804cbbac2e2")
	assert.Equal("", route)
	assert.NotNil(err)
	assert.Contains(err.Error(), "Route not found")
}
