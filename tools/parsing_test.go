package tools

import (
	asserting "github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteIdFromYaml(t *testing.T) {
	assert := asserting.New(t)
	id, err := RouteIdFromYaml(`id: 04b2e1b7-fb60-472f-a79f-af7e2353f122
type: rule_chain
attributes:
created_at: '2021-11-26T18:10:08'
destination_override_endpoint: 'https://echo.apps.verygood.systems'
entries: []`)
	assert.Nil(err)
	assert.Equal("04b2e1b7-fb60-472f-a79f-af7e2353f122", id)
}
