package leases4

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"gotest.tools/v3/assert"

	"kea-cli/api"
)

func TestGet(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	respOK := api.CommandResponse{
		Result: 0,
		Text:   "IPv4 lease found.",
		Arguments: api.Lease{
			ClientID:  "42:42:42:42:42:42:42:42",
			Cltt:      12345678,
			FQDNFwd:   false,
			FQDNRev:   true,
			Hostname:  "myhost.example.com.",
			HwAddress: "08:08:08:08:08:08",
			IpAddress: "192.0.2.1",
			State:     0,
			SubnetID:  44,
			ValidLft:  3600,
		},
	}
	bRespOK, err := json.Marshal(respOK)
	assert.NilError(t, err, "expecting nil error")

	respKO := api.CommandResponse{
		Result: 1,
		Text:   "bad parameter 'ip-address'",
	}
	bRespKO, err := json.Marshal(respKO)
	assert.NilError(t, err, "expecting nil error")

	// Start a local HTTP server simulating KEA REST API server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		assert.NilError(t, err, "expecting nil error")

		var cmd api.CommandRequest
		err = json.Unmarshal(body, &cmd)
		assert.NilError(t, err, "expecting nil error")

		if cmd.Arguments["ip-address"] == "bad" {
			rw.Write(bRespKO)
		} else {
			rw.Write(bRespOK)
		}
	}))
	defer server.Close()

	testClientConfig := api.Configuration{
		APIURL:     strings.TrimPrefix(server.URL, "http://"),
		HttpClient: server.Client(),
	}
	client, err := api.GetClient(testClientConfig)
	assert.NilError(t, err, "expecting nil error")

	type args struct {
		ipAddress string
		hwAddress string
		clientID  string
		subnetID  int
	}

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errContains string
	}{
		{"SimpleTestWithIpAddress", args{ipAddress: "192.0.2.202"}, false, ""},
		{"SimpleTestWithMACAddress", args{hwAddress: "1a:1b:1c:1d:1e:1f", clientID: "42:42:42:42:42:42:42:42"}, false, ""},
		{"SimpleTestWithClientID", args{hwAddress: "08:08:08:08:08:08:42", clientID: "42:42:42:42:42:42:42:42"}, false, ""},
		{"BadIpAddressTest", args{ipAddress: "bad"}, true, "bad parameter 'ip-address'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = get(client, tt.args.ipAddress, tt.args.hwAddress, tt.args.clientID, tt.args.subnetID)
			if tt.wantErr != (err != nil) {
				t.Errorf("get() err= %v, wantError %v", err, tt.wantErr)
			}
			if tt.wantErr == true && tt.errContains != "" {
				assert.ErrorContains(t, err, tt.errContains)
			}
		})
	}
}
