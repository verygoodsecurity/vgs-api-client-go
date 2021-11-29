package tools

import (
	asserting "github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonEquals(t *testing.T) {
	assert := asserting.New(t)

	assert.True(JsonEquals(`1`, `1`))
	assert.True(JsonEquals(`{}`, `{}`))
	assert.True(JsonEquals(`{"test": 1}`, `{"test": 1}`))
	assert.True(JsonEquals(`{"test": { "nested": [10, 11] } }`, `{"test": { "nested": [10, 11] } }`))
	assert.True(JsonEquals(`{
					"multiline": 
						{ 
					   		"some_date": "2021-11-26T18:10:08"
						} 
                                      }`, `{"multiline":{"some_date":"2021-11-26T18:10:08"}}`))

	assert.False(JsonEquals(`1`, `2`))
	assert.False(JsonEquals(`{}`, `[]`))
	assert.False(JsonEquals(`{"test": false}`, `{"test": true}`))
	assert.False(JsonEquals(`{"test": ["1", "2"]}`, `{"test": [1, 2]}`))
}

func TestYaml2Json(t *testing.T) {
	assert := asserting.New(t)
	jsonStr, err := Yaml2Json(`data:
  - attributes:
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
      updated_at: '2021-11-26T18:10:08'
    id: 04b2e1b7-fb60-472f-a79f-af7e2353f122
    type: rule_chain
`)

	assert.Nil(err)
	assert.True(JsonEquals(jsonStr, `{
"data" : [ {
  "id" : "04b2e1b7-fb60-472f-a79f-af7e2353f122",
  "type" : "rule_chain",
  "attributes" : {
    "id" : "04b2e1b7-fb60-472f-a79f-af7e2353f122",
    "created_at" : "2021-11-26T18:10:08",
    "updated_at" : "2021-11-26T18:10:08",
    "protocol" : "http",
    "source_endpoint" : "*",
    "destination_override_endpoint" : "https://echo.apps.verygood.systems",
    "host_endpoint" : "(.*)\\.verygoodproxy\\.io",
    "port" : 80,
    "ordinal" : null,
    "tags" : {
      "name" : "echo.apps.verygood.systems-beige-crescent",
      "source" : "RouteContainer"
    },
    "entries" : [ {
      "id" : "39f2f5db-06a0-461d-9387-dd9a7ab19035",
      "phase" : "REQUEST",
      "operation" : "REDACT",
      "token_manager" : "PERSISTENT",
      "public_token_generator" : "UUID",
      "transformer" : "JSON_PATH",
      "transformer_config" : [ "$.account_number" ],
      "transformer_config_map" : null,
      "operations" : null,
      "targets" : [ "body" ],
      "id_selector" : null,
      "classifiers" : { },
      "config" : {
        "condition" : "AND",
        "rules" : [ {
          "expression" : {
            "field" : "PathInfo",
            "type" : "string",
            "operator" : "matches",
            "values" : [ "/post" ]
          }
        }, {
          "rules" : null,
          "expression" : {
            "field" : "ContentType",
            "type" : "string",
            "operator" : "equals",
            "values" : [ "application/json" ]
          }
        } ]
      }
    } ]
  }
} ]
	}`))

}
