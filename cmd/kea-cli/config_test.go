package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"gotest.tools/v3/assert"

	"kea-cli/api"
)

func TestDefaultClientConfig(t *testing.T) {
	initFlags()

	expected := &api.Configuration{
		APIURL:            api.DefaultKeaAPIURL,
		HTTPClientTimeout: api.DefaultHTTPClientTimeout,
		SSLEnabled:        false,
		SkipTLSVerify:     false,
		LogLevel:          log.InfoLevel.String(),
	}
	actual := new(api.Configuration)

	err := getClientConfig(actual)
	assert.NilError(t, err, "expecting nil error")
	assert.DeepEqual(t, expected, actual)
}
